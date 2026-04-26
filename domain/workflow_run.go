package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkflowRun struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Status WorkflowRunStatus `gorm:"type:varchar(20);not null;check:status IN ('pending','running','completed','failed')" json:"status"`

	WorkflowID uuid.UUID `gorm:"type:uuid;not null" json:"workflow_id"`
	Workflow   Workflow  `gorm:"foreignKey:WorkflowID" json:"workflow"`

	StartedAt   *time.Time `gorm:"null" json:"started_at"`
	CompletedAt *time.Time `gorm:"null" json:"completed_at"`
	FailedAt    *time.Time `gorm:"null" json:"failed_at"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (WorkflowRun) TableName() string {
	return "workflow_runs"
}
