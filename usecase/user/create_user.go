package user

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
)

type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, clerkID string, firstName string, lastName string, banned bool) (*entity.User, error) {
	user := entity.User{
		ClerkID:   clerkID,
		FirstName: firstName,
		LastName:  lastName,
		Banned:    banned,
	}

	err := u.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}
