package entity

import (
	"time"

	"flowforge-api/domain/types"

	"github.com/google/uuid"
)

type Endpoint struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name    string    `gorm:"not null;index:idx_endpoint_name" json:"name"`
	BaseURI string    `gorm:"not null" json:"baseUri"`
	Path    string    `gorm:"not null" json:"path"`
	Method  string    `gorm:"not null" json:"method"`
	Timeout int       `gorm:"not null;default:30" json:"timeout"`

	RetryOnFailure bool `gorm:"not null;default:false" json:"retryOnFailure"`
	RetryCount     int  `gorm:"not null;default:0" json:"retryCount"`
	RetryDelay     int  `gorm:"not null;default:0" json:"retryDelay"`

	Query  types.Query  `json:"query" gorm:"type:jsonb;default:'[]'"`
	Header types.Header `json:"header" gorm:"type:jsonb;default:'[]'"`
	Body   types.Body   `json:"body" gorm:"type:jsonb;default:'[]'"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_endpoint_org" json:"organization_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Endpoint) TableName() string {
	return "endpoints"
}
