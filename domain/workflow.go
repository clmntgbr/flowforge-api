package domain

import (
	"time"

	"github.com/google/uuid"
)

type Workflow struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"null" json:"description"`

	Steps          []Step      `gorm:"foreignKey:WorkflowID" json:"steps"`
	Connexions     []Connexion `gorm:"foreignKey:WorkflowID" json:"connexions"`
	OrganizationID uuid.UUID   `gorm:"type:uuid;not null" json:"organization_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Workflow) TableName() string {
	return "workflows"
}
