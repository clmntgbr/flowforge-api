package consumer

import (
	"encoding/json"

	"flowforge-api/infrastructure/runner"
)

type ConsumerMessage struct {
	SecretKey string          `json:"secret_key" validate:"required"`
	Message   json.RawMessage `json:"message" validate:"required"`
}

type ConsumerCompletedMessage struct {
	WorkflowRunID string                `json:"workflow_run_id" validate:"required,uuid"`
	StepRunID     string                `json:"step_run_id" validate:"required,uuid"`
	CompletedAt   string                `json:"completed_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Insights      runner.RunnerInsights `json:"insights" validate:"required"`
	Response      string                `json:"response" validate:"required"`
}

type ConsumerFailedMessage struct {
	WorkflowRunID string                `json:"workflow_run_id" validate:"required,uuid"`
	StepRunID     string                `json:"step_run_id" validate:"required,uuid"`
	Error         string                `json:"error" validate:"required,max=2048"`
	FailedAt      string                `json:"failed_at" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Insights      runner.RunnerInsights `json:"insights" validate:"required"`
	Response      string                `json:"response"`
}
