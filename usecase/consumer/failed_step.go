package consumer

import (
	"context"
	"errors"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/usecase/insight"
	"time"

	"github.com/google/uuid"
)

type FailedStepUseCase struct {
	createInsightUseCase *insight.CreateInsightUseCase
	stepRunRepo          *repository.StepRunRepository
	workflowRunRepo      *repository.WorkflowRunRepository
}

func NewFailedStepUseCase(createInsightUseCase *insight.CreateInsightUseCase, stepRunRepo *repository.StepRunRepository, workflowRunRepo *repository.WorkflowRunRepository) *FailedStepUseCase {
	return &FailedStepUseCase{
		createInsightUseCase: createInsightUseCase,
		stepRunRepo:          stepRunRepo,
		workflowRunRepo:      workflowRunRepo,
	}
}

func (u *FailedStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerFailedMessage) error {
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

	failedAt, err := time.Parse(time.RFC3339, message.FailedAt)
	if err != nil {
		return err
	}

	stepRun.Status = enum.StepRunStatusFailed
	stepRun.FailedAt = &failedAt
	stepRun.Error = message.Error
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

	workflowRun.FailedAt = &failedAt
	workflowRun.Status = enum.WorkflowRunStatusFailed
	workflowRun.Error = stepRun.Error
	workflowRun.ExecutedSteps = append(workflowRun.ExecutedSteps, stepRun.StepID.String())

	err = (*u.workflowRunRepo).Update(ctx, workflowRun)
	if err != nil {
		return err
	}

	return nil
}
