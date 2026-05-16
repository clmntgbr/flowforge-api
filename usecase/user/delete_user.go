package user

import (
	"context"
	"errors"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type DeleteUserUseCase struct {
	userRepo *repository.UserRepository
}

func NewDeleteUserUseCase(userRepo *repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{userRepo: userRepo}
}

func (u *DeleteUserUseCase) Execute(ctx context.Context, userID string) error {
	err := (*u.userRepo).Delete(ctx, uuid.MustParse(userID))
	if err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}
