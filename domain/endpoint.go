package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Endpoint struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name    string    `gorm:"not null" json:"name"`
	BaseURI string    `gorm:"not null" json:"baseUri"`
	Path    string    `gorm:"not null" json:"path"`
	Method  string    `gorm:"not null" json:"method"`
	Timeout int       `gorm:"not null;default:30" json:"timeout"`

	Query datatypes.JSON `json:"query" gorm:"type:jsonb;default:'[]'"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null" json:"organization_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Endpoint) TableName() string {
	return "endpoints"
}
