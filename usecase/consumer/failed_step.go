package consumer

import (
	"context"
	"errors"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/config"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/usecase/insight"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow_run"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type FailedStepUseCase struct {
	createInsightUseCase         *insight.CreateInsightUseCase
	stepRunRepo                  *repository.StepRunRepository
	workflowRunRepo              *repository.WorkflowRunRepository
	stepRepo                     *repository.StepRepository
	createStepRunUseCase         *step_run.CreateStepRunUseCase
	executeStepRunUseCase        *step_run.ExecuteStepRunUseCase
	isCanceledWorkflowRunUseCase *workflow_run.IsCanceledWorkflowRunUseCase
	stepRunPublisher             rabbitmq.Publisher
	env                          *config.Config
}

func NewFailedStepUseCase(
	createInsightUseCase *insight.CreateInsightUseCase,
	stepRunRepo *repository.StepRunRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
	stepRepo *repository.StepRepository,
	createStepRunUseCase *step_run.CreateStepRunUseCase,
	executeStepRunUseCase *step_run.ExecuteStepRunUseCase,
	isCanceledWorkflowRunUseCase *workflow_run.IsCanceledWorkflowRunUseCase,
	stepRunPublisher rabbitmq.Publisher,
	env *config.Config,
) *FailedStepUseCase {
	return &FailedStepUseCase{
		createInsightUseCase:         createInsightUseCase,
		stepRunRepo:                  stepRunRepo,
		workflowRunRepo:              workflowRunRepo,
		stepRepo:                     stepRepo,
		createStepRunUseCase:         createStepRunUseCase,
		executeStepRunUseCase:        executeStepRunUseCase,
		isCanceledWorkflowRunUseCase: isCanceledWorkflowRunUseCase,
		stepRunPublisher:             stepRunPublisher,
		env:                          env,
	}
}

func (u *FailedStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerFailedMessage) error {
	err := u.isCanceledWorkflowRunUseCase.Execute(ctx, message.WorkflowRunID)
	if err != nil {
		return nil
	}

	ins, err := u.createInsightUseCase.Execute(
		ctx,
		message.Insights.StartTime,
		message.Insights.EndTime,
		message.Insights.Duration,
		message.Insights.StatusCode,
		message.Insights.ResponseSize,
		message.Insights.AttemptNumber,
		message.Insights.TotalAttempts,
		message.Insights.QueueTime,
		message.Insights.DNSLookupDuration,
		message.Insights.TCPConnectionTime,
		message.Insights.TLSHandshakeTime,
		message.Insights.TTFB,
		message.Insights.ErrorMessage,
		message.Insights.ErrorType,
		message.Insights.RequestSize,
	)
	if err != nil {
		return err
	}

	stepRunID := uuid.MustParse(message.StepRunID)

	stepRun, err := (*u.stepRunRepo).GetByID(ctx, stepRunID)
	if err != nil {
		return err
	}
	if stepRun == nil {
		return errors.New("step run not found")
	}

	failedAt := message.FailedAt.UTC()
	stepRun.Status = enum.StepRunStatusFailed
	stepRun.Statuses = append(stepRun.Statuses, enum.StepRunStatusFailed)
	stepRun.FailedAt = &failedAt
	stepRun.CompletedAt = &failedAt
	stepRun.Error = message.Error
	stepRun.Response = message.Response
	stepRun.InsightID = &ins.ID

	if err = (*u.stepRunRepo).Update(ctx, stepRun); err != nil {
		return err
	}

	workflowRun, err := (*u.workflowRunRepo).GetByID(ctx, stepRun.WorkflowRunID)
	if err != nil {
		return err
	}
	if workflowRun == nil {
		return errors.New("workflow run not found")
	}

	workflowRun.ExecutedSteps = append(workflowRun.ExecutedSteps, stepRun.StepID.String())
	workflowRun.FailedSteps = append(workflowRun.FailedSteps, stepRun.StepID.String())

	failedStep, err := (*u.stepRepo).GetByID(ctx, stepRun.StepID)
	if err != nil {
		return fmt.Errorf("failed to get failed step: %w", err)
	}

	nextCandidate, err := (*u.stepRepo).GetFirstStepAtLevel(ctx, workflowRun.WorkflowID, failedStep.TreeIndex, majorFromIndex(failedStep.Index), failedStep.ExecutionOrder, workflowRun.ExecutedSteps)
	if err != nil {
		return fmt.Errorf("failed to find alternative step: %w", err)
	}

	if nextCandidate == nil {
		nextCandidate, err = (*u.stepRepo).GetNextStepByWorkflowID(ctx, workflowRun.WorkflowID, failedStep.TreeIndex, workflowRun.ExecutedSteps)
		if err != nil {
			return fmt.Errorf("failed to find next tree step: %w", err)
		}
	}

	if nextCandidate == nil {
		workflowRun.FailedAt = &failedAt
		workflowRun.CompletedAt = &failedAt
		workflowRun.Status = enum.WorkflowRunStatusFailed
		workflowRun.Statuses = append(workflowRun.Statuses, enum.WorkflowRunStatusFailed)

		return (*u.workflowRunRepo).Update(ctx, workflowRun)
	}

	if err = (*u.workflowRunRepo).Update(ctx, workflowRun); err != nil {
		return err
	}

	nextStepRun, err := u.createStepRunUseCase.Execute(ctx, workflowRun.ID, nextCandidate.ID)
	if err != nil {
		return fmt.Errorf("failed to create alternative step run: %w", err)
	}

	nextStepRun, err = u.executeStepRunUseCase.Execute(ctx, &nextStepRun)
	if err != nil {
		return fmt.Errorf("failed to execute alternative step run: %w", err)
	}

	event := rabbitmq.NewStepRunEvent(nextStepRun)
	if err := u.stepRunPublisher.PublishStepRunEvent(ctx, u.env, event); err != nil {
		return fmt.Errorf("failed to publish alternative step run: %w", err)
	}

	return nil
}

func majorFromIndex(index string) int {
	parts := strings.Split(index, ".")
	v, _ := strconv.Atoi(parts[0])
	return v
}
