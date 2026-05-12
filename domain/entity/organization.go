package entity

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"not null" json:"name"`

	Users []User `gorm:"many2many:user_organizations" json:"users,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	IsActive bool `gorm:"-" json:"isActive"`
}

func (Organization) TableName() string {
	return "organizations"
}
