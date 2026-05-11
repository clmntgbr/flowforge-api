package user

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
)

type UpdateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: userRepo}
}

func (s *UpdateUserUseCase) Execute(ctx context.Context, user *entity.User) error {
	err := s.userRepo.Update(ctx, user)
	if err != nil {
		return errors.New("failed to update user")
	}

	return nil
}
