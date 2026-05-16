package consumer

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/usecase/insight"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CompletedStepUseCase struct {
	createInsightUseCase insight.CreateInsightUseCase
	stepRunRepo          repository.StepRunRepository
	workflowRunRepo      repository.WorkflowRunRepository
	stepRepo             repository.StepRepository
}

func NewCompletedStepUseCase(createInsightUseCase insight.CreateInsightUseCase, stepRunRepo repository.StepRunRepository, workflowRunRepo repository.WorkflowRunRepository, stepRepo repository.StepRepository) *CompletedStepUseCase {
	return &CompletedStepUseCase{
		createInsightUseCase: createInsightUseCase,
		stepRunRepo:          stepRunRepo,
		workflowRunRepo:      workflowRunRepo,
		stepRepo:             stepRepo,
	}
}

func (u *CompletedStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerCompletedMessage) error {
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

	stepRun, err := u.stepRunRepo.GetByID(ctx, stepRunID)
	if err != nil {
		return err
	}

	if stepRun == nil {
		return errors.New("step run not found")
	}

	completedAt, err := time.Parse(time.RFC3339, message.CompletedAt)
	if err != nil {
		return err
	}

	stepRun.Status = enum.StepRunStatusCompleted
	stepRun.CompletedAt = &completedAt
	stepRun.Response = message.Response
	stepRun.InsightID = &insight.ID

	err = u.stepRunRepo.Update(ctx, stepRun)
	if err != nil {
		return err
	}

	workflowRun, err := u.workflowRunRepo.GetByID(ctx, stepRun.WorkflowRunID)
	if err != nil {
		return err
	}

	if workflowRun == nil {
		return errors.New("workflow run not found")
	}

	workflowRun.ExecutedSteps = append(workflowRun.ExecutedSteps, stepRun.StepID.String())

	err = u.workflowRunRepo.Update(ctx, workflowRun)
	if err != nil {
		return err
	}

	nextStep, err := u.stepRepo.GetNextStepByWorkflowID(ctx, workflowRun.WorkflowID, workflowRun.ExecutedSteps)
	if err != nil {
		return err
	}

	if nextStep == nil {
		return errors.New("no next step found")
	}

	if nextStep == nil {
		err = u.completeWorkflowRun(ctx, workflowRun, message)
		if err != nil {
			return fmt.Errorf("failed to complete workflow run: %w", err)
		}
		return nil
	}

	return nil
}

func (u *CompletedStepUseCase) completeWorkflowRun(ctx context.Context, workflowRun *entity.WorkflowRun, message consumerDTO.ConsumerCompletedMessage) error {
	completedAt, err := time.Parse(time.RFC3339, message.CompletedAt)
	if err != nil {
		return err
	}

	workflowRun.CompletedAt = &completedAt
	workflowRun.Status = enum.WorkflowRunStatusCompleted

	err = u.workflowRunRepo.Update(ctx, workflowRun)
	if err != nil {
		return err
	}

	return nil
}
