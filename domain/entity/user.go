package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ClerkID   string    `gorm:"uniqueIndex;not null" json:"clerk_id"`
	FirstName string    `gorm:"null" json:"first_name"`
	LastName  string    `gorm:"null" json:"last_name"`
	Banned    bool      `gorm:"default:false;index:idx_user_banned" json:"banned"`

	Organizations        []Organization `gorm:"many2many:user_organizations" json:"organizations,omitempty"`
	ActiveOrganization   *Organization  `gorm:"foreignKey:ActiveOrganizationID" json:"active_organization,omitempty"`
	ActiveOrganizationID *uuid.UUID     `gorm:"type:uuid;default:null;index:idx_user_active_org" json:"active_organization_id,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
