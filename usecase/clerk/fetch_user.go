package clerk

import (
	"context"
	"errors"
	"flowforge-api/infrastructure/config"
	"flowforge-api/presenter"

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

func (s *FetchUserUseCase) Execute(ctx context.Context, clerkID string) (presenter.ClerkUser, error) {
	clerkUser, err := clerkuser.Get(context.Background(), clerkID)
	if err != nil {
		return presenter.ClerkUser{}, errors.New("failed to get user")
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

	return presenter.ClerkUser{
		ID:        clerkUser.ID,
		FirstName: firstName,
		LastName:  lastName,
		Banned:    banned,
	}, nil
}
