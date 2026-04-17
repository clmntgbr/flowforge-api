package domain

import (
	"time"

	"github.com/google/uuid"
)

type Connexion struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	FromStepID uuid.UUID `gorm:"type:uuid;not null" json:"from_step_id"`
	ToStepID   uuid.UUID `gorm:"type:uuid;not null" json:"to_step_id"`
	WorkflowID uuid.UUID `gorm:"type:uuid;not null" json:"workflow_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Connexion) TableName() string {
	return "connexions"
}
