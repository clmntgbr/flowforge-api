package clerk

import (
	"context"
	"errors"
	clerkdto "flowforge-api/infrastructure/clerk"
	"flowforge-api/infrastructure/config"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
)

type FetchUserUseCase struct {
	config *config.Config
}

func NewFetchUserUseCase(cfg *config.Config) *FetchUserUseCase {
	clerk.SetKey(cfg.ClerkSecretKey)
	return &FetchUserUseCase{config: cfg}
}

func (s *FetchUserUseCase) Execute(ctx context.Context, clerkID string) (clerkdto.ClerkUser, error) {
	clerkUser, err := clerkuser.Get(context.Background(), clerkID)
	if err != nil {
		return clerkdto.ClerkUser{}, errors.New("failed to get user")
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

	return clerkdto.ClerkUser{
		ID:        clerkUser.ID,
		FirstName: firstName,
		LastName:  lastName,
		Banned:    banned,
	}, nil
}
