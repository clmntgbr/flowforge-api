package domain

import (
	"time"

	"github.com/google/uuid"
)

type StepRun struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Status StepRunStatus `gorm:"type:varchar(20);not null;check:status IN ('pending','running','completed','failed')" json:"status"`

	StepID uuid.UUID `gorm:"type:uuid;not null" json:"step_id"`
	Step   Step      `gorm:"foreignKey:StepID" json:"step"`

	StartedAt   *time.Time `gorm:"null" json:"started_at"`
	CompletedAt *time.Time `gorm:"null" json:"completed_at"`
	FailedAt    *time.Time `gorm:"null" json:"failed_at"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (StepRun) TableName() string {
	return "step_runs"
}
