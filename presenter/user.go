package presenter

import (
	"flowforge-api/domain/entity"
	"time"
)

type UserDetailResponse struct {
	ID        string    `json:"id"`
	ClerkID   string    `json:"clerkId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUserDetailResponse(user entity.User) UserDetailResponse {
	return UserDetailResponse{
		ID:        user.ID.String(),
		ClerkID:   user.ClerkID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
