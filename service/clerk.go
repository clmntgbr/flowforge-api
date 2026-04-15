package service

import (
	"context"
	"go-api/config"
	"go-api/dto"
	"go-api/errors"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
)

type ClerkService struct {
	config *config.Config
}

func NewClerkService(cfg *config.Config) *ClerkService {
	clerk.SetKey(cfg.ClerkSecretKey)
	return &ClerkService{config: cfg}
}

func (s *ClerkService) GetUser(clerkID string) (dto.ClerkUser, error) {
	clerkUser, err := clerkuser.Get(context.Background(), clerkID)
	if err != nil {
		return dto.ClerkUser{}, errors.ErrUserNotFound
	}

	firstName := ""
	if clerkUser.FirstName != nil {
		firstName = *clerkUser.FirstName
	}

	lastName := ""
	if clerkUser.LastName != nil {
		lastName = *clerkUser.LastName
	}

	banned := clerkUser.Banned

	return dto.ClerkUser{
		ID:        clerkUser.ID,
		FirstName: firstName,
		LastName:  lastName,
		Banned:    banned,
	}, nil
}
