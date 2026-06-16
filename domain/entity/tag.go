package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name  string    `gorm:"not null;index:idx_endpoint_name" json:"name"`
	Color string    `gorm:"not null" json:"color"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Tag) TableName() string {
	return "tags"
}
