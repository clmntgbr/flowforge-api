package rabbitmq

import (
	"flowforge-api/domain/entity"
	"flowforge-api/infrastructure/runner"
	"flowforge-api/presenter"
	"time"

	"github.com/google/uuid"
)

type MessagePayload struct {
	SecretKey    string       `json:"secret_key"`
	StepRunEvent StepRunEvent `json:"step_run_event"`
}

type MessageResponse struct {
	SecretKey string                         `json:"secret_key" validate:"required"`
	Message   RunnerMessageResponseInterface `json:"message" validate:"required"`
}

type StepRunEvent struct {
	WorkflowRunID uuid.UUID                        `json:"workflow_run_id" validate:"required,uuid"`
	StepRunID     uuid.UUID                        `json:"step_run_id" validate:"required,uuid"`
	WorkflowID    uuid.UUID                        `json:"workflow_id" validate:"required,uuid"`
	Step          presenter.StepDetailResponse     `json:"step" validate:"required"`
	Endpoint      presenter.EndpointDetailResponse `json:"endpoint" validate:"required"`
	QueuedAt      *time.Time                       `json:"queued_at,omitempty" validate:"omitempty"`
}

type RunnerMessageResponseInterface interface {
	isMessage()
}

func (RunnerCompletedMessage) isMessage() {}
func (RunnerFailedMessage) isMessage()    {}

type RunnerCompletedMessage struct {
	WorkflowRunID string                `json:"workflow_run_id" validate:"required,uuid"`
	StepRunID     string                `json:"step_run_id" validate:"required,uuid"`
	CompletedAt   string                `json:"completed_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Insights      runner.RunnerInsights `json:"insights" validate:"required"`
	Response      string                `json:"response" validate:"required"`
}

type RunnerFailedMessage struct {
	WorkflowRunID string                `json:"workflow_run_id" validate:"required,uuid"`
	StepRunID     string                `json:"step_run_id" validate:"required,uuid"`
	Error         string                `json:"error" validate:"required,max=2048"`
	FailedAt      string                `json:"failed_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Insights      runner.RunnerInsights `json:"insights" validate:"required"`
	Response      string                `json:"response"`
}

func NewStepRunEvent(stepRun entity.StepRun) StepRunEvent {
	now := time.Now().UTC()

	return StepRunEvent{
		WorkflowRunID: stepRun.WorkflowRunID,
		StepRunID:     stepRun.ID,
		WorkflowID:    stepRun.Step.WorkflowID,
		Step:          presenter.NewStepDetailResponse(stepRun.Step),
		Endpoint:      presenter.NewEndpointDetailResponse(stepRun.Step.Endpoint),
		QueuedAt:      &now,
	}
}
