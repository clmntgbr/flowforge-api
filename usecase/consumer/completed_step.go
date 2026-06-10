package consumer

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/config"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/usecase/insight"
	usecaseStep "flowforge-api/usecase/step"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow_run"
	"fmt"

	"github.com/google/uuid"
)

type CompletedStepUseCase struct {
	createInsightUseCase         *insight.CreateInsightUseCase
	stepRunRepo                  *repository.StepRunRepository
	workflowRunRepo              *repository.WorkflowRunRepository
	findNextStepUseCase          *usecaseStep.FindNextStepUseCase
	createStepRunUseCase         *step_run.CreateStepRunUseCase
	executeStepRunUseCase        *step_run.ExecuteStepRunUseCase
	isCanceledWorkflowRunUseCase *workflow_run.IsCanceledWorkflowRunUseCase
	stepRunPublisher             rabbitmq.Publisher
	env                          *config.Config
}

func NewCompletedStepUseCase(
	createInsightUseCase *insight.CreateInsightUseCase,
	stepRunRepo *repository.StepRunRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
	findNextStepUseCase *usecaseStep.FindNextStepUseCase,
	createStepRunUseCase *step_run.CreateStepRunUseCase,
	executeStepRunUseCase *step_run.ExecuteStepRunUseCase,
	isCanceledWorkflowRunUseCase *workflow_run.IsCanceledWorkflowRunUseCase,
	stepRunPublisher rabbitmq.Publisher,
	env *config.Config,
) *CompletedStepUseCase {
	return &CompletedStepUseCase{
		createInsightUseCase:         createInsightUseCase,
		stepRunRepo:                  stepRunRepo,
		workflowRunRepo:              workflowRunRepo,
		findNextStepUseCase:          findNextStepUseCase,
		createStepRunUseCase:         createStepRunUseCase,
		executeStepRunUseCase:        executeStepRunUseCase,
		isCanceledWorkflowRunUseCase: isCanceledWorkflowRunUseCase,
		stepRunPublisher:             stepRunPublisher,
		env:                          env,
	}
}

func (u *CompletedStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerCompletedMessage) error {
	err := u.isCanceledWorkflowRunUseCase.Execute(ctx, message.WorkflowRunID)
	if err != nil {
		return nil
	}

	insight, err := u.createInsightUseCase.Execute(
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

	completedAt := message.CompletedAt.UTC()
	stepRun.Status = enum.StepRunStatusCompleted
	stepRun.Statuses = append(stepRun.Statuses, enum.StepRunStatusCompleted)
	stepRun.CompletedAt = &completedAt
	stepRun.Response = message.Response
	stepRun.InsightID = &insight.ID

	err = (*u.stepRunRepo).Update(ctx, stepRun)
	if err != nil {
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

	err = (*u.workflowRunRepo).Update(ctx, workflowRun)
	if err != nil {
		return err
	}

	nextStep, err := u.findNextStepUseCase.Execute(ctx, workflowRun.WorkflowID, workflowRun.ExecutedSteps)
	if err != nil {
		return err
	}

	if nextStep == nil {
		err = u.completeWorkflowRun(ctx, workflowRun, message)
		if err != nil {
			return fmt.Errorf("failed to complete workflow run: %w", err)
		}
		return nil
	}

	nextStepRun, err := u.createStepRunUseCase.Execute(ctx, workflowRun.ID, nextStep.ID)
	if err != nil {
		return err
	}

	nextStepRun, err = u.executeStepRunUseCase.Execute(ctx, &nextStepRun)
	if err != nil {
		return err
	}

	event := rabbitmq.NewStepRunEvent(nextStepRun)
	if err := u.stepRunPublisher.PublishStepRunEvent(ctx, u.env, event); err != nil {
		return fmt.Errorf("🚨 failed to publish step run: %w", err)
	}

	return nil
}

func (u *CompletedStepUseCase) completeWorkflowRun(ctx context.Context, workflowRun *entity.WorkflowRun, message consumerDTO.ConsumerCompletedMessage) error {
	completedAt := message.CompletedAt.UTC()
	workflowRun.CompletedAt = &completedAt
	workflowRun.Status = enum.WorkflowRunStatusCompleted
	workflowRun.Statuses = append(workflowRun.Statuses, enum.WorkflowRunStatusCompleted)

	err := (*u.workflowRunRepo).Update(ctx, workflowRun)
	if err != nil {
		return err
	}

	return nil
}
