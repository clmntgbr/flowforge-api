package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ClerkID   string    `gorm:"uniqueIndex;not null" json:"clerk_id"`
	FirstName string    `gorm:"null" json:"first_name"`
	LastName  string    `gorm:"null" json:"last_name"`
	Banned    bool      `gorm:"default:false" json:"banned"`

	Projects        []Project `gorm:"many2many:user_projects" json:"projects,omitempty"`
	ActiveProject   Project   `gorm:"foreignKey:ActiveProjectID" json:"active_project,omitempty"`
	ActiveProjectID uuid.UUID `gorm:"type:uuid" json:"active_project_id"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
