package dto

import (
	"go-api/domain"
	"time"
)

type UserOutput struct {
	ID        string    `json:"id"`
	ClerkID   string    `json:"clerkId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUserOutput(user domain.User) UserOutput {
	return UserOutput{
		ID:        user.ID.String(),
		ClerkID:   user.ClerkID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
