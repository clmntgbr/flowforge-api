package entity

import (
	"time"

	"github.com/google/uuid"
)

type Connexion struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	FromStepID uuid.UUID `gorm:"type:uuid;not null;index:idx_connexion_from_to,priority:1;index:idx_connexion_from" json:"from_step_id"`
	ToStepID   uuid.UUID `gorm:"type:uuid;not null;index:idx_connexion_from_to,priority:2;index:idx_connexion_to" json:"to_step_id"`
	WorkflowID uuid.UUID `gorm:"type:uuid;not null;index:idx_connexion_workflow" json:"workflow_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Connexion) TableName() string {
	return "connexions"
}
