package rabbitmq

import (
	"flowforge-api/domain/entity"
	"flowforge-api/presenter"
	"time"

	"github.com/google/uuid"
)

type MessagePayload struct {
	SecretKey    string       `json:"secret_key"`
	StepRunEvent StepRunEvent `json:"step_run_event"`
}

type StepRunEvent struct {
	WorkflowRunID uuid.UUID                        `json:"workflow_run_id" validate:"required,uuid"`
	StepRunID     uuid.UUID                        `json:"step_run_id" validate:"required,uuid"`
	WorkflowID    uuid.UUID                        `json:"workflow_id" validate:"required,uuid"`
	Step          presenter.StepDetailResponse     `json:"step" validate:"required"`
	Endpoint      presenter.EndpointDetailResponse `json:"endpoint" validate:"required"`
	QueuedAt      *time.Time                       `json:"queued_at,omitempty" validate:"omitempty"`
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
