package entity

import (
	"time"

	"github.com/google/uuid"
)

type Variable struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:varchar(255);null" json:"description"`
	Path        string    `gorm:"type:varchar(255);not null" json:"path"`
	StepID      uuid.UUID `gorm:"type:uuid;not null;index:idx_variable_step" json:"step_id"`
	WorkflowID  uuid.UUID `gorm:"type:uuid;not null;index:idx_variable_workflow" json:"workflow_id"`

	IsSecret     bool   `gorm:"default:false" json:"is_secret"`
	DefaultValue string `gorm:"null" json:"default_value"`
	LastValue    string `gorm:"null" json:"last_value"`

	Step     Step     `gorm:"foreignKey:StepID" json:"step"`
	Workflow Workflow `gorm:"foreignKey:WorkflowID" json:"workflow"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Variable) TableName() string {
	return "variables"
}
