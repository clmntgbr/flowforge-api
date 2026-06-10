package consumer

import (
	"encoding/json"
	"time"

	"flowforge-api/infrastructure/runner"
)

type ConsumerMessage struct {
	SecretKey string          `json:"secret_key" validate:"required"`
	Message   json.RawMessage `json:"message" validate:"required"`
}

type ConsumerCompletedMessage struct {
	WorkflowRunID string                `json:"workflow_run_id" validate:"required,uuid"`
	WorkflowID    string                `json:"workflow_id" validate:"required,uuid"`
	StepRunID     string                `json:"step_run_id" validate:"required,uuid"`
	CompletedAt   time.Time             `json:"completed_at" validate:"required"`
	Insights      runner.RunnerInsights `json:"insights" validate:"required"`
	Response      string                `json:"response" validate:"required"`
}

type ConsumerFailedMessage struct {
	WorkflowRunID string                `json:"workflow_run_id" validate:"required,uuid"`
	WorkflowID    string                `json:"workflow_id" validate:"required,uuid"`
	StepRunID     string                `json:"step_run_id" validate:"required,uuid"`
	Error         string                `json:"error" validate:"required,max=2048"`
	FailedAt      time.Time             `json:"failed_at" validate:"required"`
	Insights      runner.RunnerInsights `json:"insights" validate:"required"`
	Response      string                `json:"response"`
}
