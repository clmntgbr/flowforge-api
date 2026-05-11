package entity

import (
	"time"

	"flowforge-api/domain/enum"

	"github.com/google/uuid"
)

type Workflow struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null;index:idx_workflow_name" json:"name"`
	Description string    `gorm:"null" json:"description"`

	Status enum.WorkflowStatus `gorm:"type:varchar(20);not null;check:status IN ('active','inactive','deleted');default:inactive;index:idx_workflow_status;index:idx_workflow_org_status,priority:2" json:"status"`

	Steps          []Step      `gorm:"foreignKey:WorkflowID" json:"steps"`
	Connexions     []Connexion `gorm:"foreignKey:WorkflowID" json:"connexions"`
	OrganizationID uuid.UUID   `gorm:"type:uuid;not null;index:idx_workflow_org;index:idx_workflow_org_status,priority:1" json:"organization_id"`

	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_workflow_created" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Workflow) TableName() string {
	return "workflows"
}
