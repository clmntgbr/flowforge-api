package entity

import (
	"flowforge-api/domain/enum"
	"time"

	"github.com/google/uuid"
)

type WorkflowRun struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Status enum.WorkflowRunStatus `gorm:"type:varchar(20);not null;check:status IN ('pending','running','completed','failed');index:idx_workflow_run_status;index:idx_workflow_run_workflow_status,priority:2" json:"status"`

	WorkflowID uuid.UUID `gorm:"type:uuid;not null;index:idx_workflow_run_workflow;index:idx_workflow_run_workflow_status,priority:1" json:"workflow_id"`
	Workflow   Workflow  `gorm:"foreignKey:WorkflowID" json:"workflow"`

	ExecutedSteps []string `gorm:"serializer:json;type:jsonb;default:'[]'" json:"executed_steps"`

	StepsRuns []StepRun `gorm:"foreignKey:WorkflowRunID" json:"steps_runs"`

	StartedAt   *time.Time `gorm:"null" json:"started_at"`
	CompletedAt *time.Time `gorm:"null" json:"completed_at"`
	FailedAt    *time.Time `gorm:"null" json:"failed_at"`

	Error string `gorm:"null" json:"error"`

	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_workflow_run_created" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (WorkflowRun) TableName() string {
	return "workflow_runs"
}
