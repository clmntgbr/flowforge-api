package service

import (
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/repository"
	"time"

	"github.com/gofiber/fiber/v3"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(c fiber.Ctx, id string, firstName string, lastName string, banned bool) (*domain.User, error) {
	user := &domain.User{
		ClerkID:   id,
		FirstName: firstName,
		LastName:  lastName,
		Banned:    banned,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(c fiber.Ctx, id string, firstName string, lastName string, banned bool) error {
	user, err := s.userRepository.FindByClerkID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Banned = banned

	return s.userRepository.Update(user)
}

func (s *UserService) DeleteUser(c fiber.Ctx, id string) error {
	user, err := s.userRepository.FindByClerkID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	return s.userRepository.Delete(user)
}

func (s *UserService) GetUser(user *domain.User) (*dto.UserOutput, error) {
	output := dto.NewUserOutput(*user)
	return &output, nil
}
