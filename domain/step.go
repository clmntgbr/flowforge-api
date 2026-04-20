package domain

import (
	"time"

	"github.com/google/uuid"
)

type Position struct {
	X int `gorm:"not null" json:"x"`
	Y int `gorm:"not null" json:"y"`
}

type Step struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:varchar(255);null" json:"description"`
	Position    Position  `gorm:"embedded;embeddedPrefix:position_" json:"position"`
	Index       string    `gorm:"type:varchar(255);not null" json:"index"`
	Timeout     int       `gorm:"type:int;null" json:"timeout"`

	EndpointID uuid.UUID `gorm:"type:uuid;not null" json:"endpoint_id"`
	WorkflowID uuid.UUID `gorm:"type:uuid;not null" json:"workflow_id"`

	Endpoint Endpoint `gorm:"foreignKey:EndpointID" json:"endpoint"`

	Query  Query  `json:"query" gorm:"type:jsonb;default:'[]'"`
	Header Header `json:"headers" gorm:"type:jsonb;default:'[]'"`
	Body   Body   `json:"body" gorm:"type:jsonb;default:'[]'"`

	RetryOnFailure bool `gorm:"not null;default:false" json:"retryOnFailure"`
	RetryCount     int  `gorm:"not null;default:0" json:"retryCount"`
	RetryDelay     int  `gorm:"not null;default:0" json:"retryDelay"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Step) TableName() string {
	return "steps"
}
