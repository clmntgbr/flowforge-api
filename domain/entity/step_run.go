package entity

import (
	"flowforge-api/domain/enum"
	"time"

	"github.com/google/uuid"
)

type StepRun struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Status   enum.StepRunStatus   `gorm:"type:varchar(20);not null;check:status IN ('pending','running','completed','failed');index:idx_step_run_status;index:idx_step_run_workflow_status,priority:2" json:"status"`
	Statuses []enum.StepRunStatus `gorm:"serializer:json;type:jsonb;default:'[]'" json:"statuses"`

	StepID uuid.UUID `gorm:"type:uuid;not null;index:idx_step_run_step" json:"step_id"`
	Step   Step      `gorm:"foreignKey:StepID" json:"step"`

	WorkflowRunID uuid.UUID   `gorm:"type:uuid;not null;index:idx_step_run_workflow_run;index:idx_step_run_workflow_status,priority:1" json:"workflow_run_id"`
	WorkflowRun   WorkflowRun `gorm:"foreignKey:WorkflowRunID" json:"workflow_run"`

	StartedAt   *time.Time `gorm:"null" json:"started_at"`
	CompletedAt *time.Time `gorm:"null" json:"completed_at"`
	FailedAt    *time.Time `gorm:"null" json:"failed_at"`

	Error    string `gorm:"null" json:"error"`
	Response string `gorm:"null" json:"response"`

	InsightID *uuid.UUID `gorm:"type:uuid;index:idx_step_run_insight" json:"insight_id"`
	Insight   *Insight   `gorm:"foreignKey:InsightID" json:"insight"`

	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_step_run_created" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StepRun) TableName() string {
	return "step_runs"
}
