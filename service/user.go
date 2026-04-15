package service

import (
	"go-api/domain"
	"go-api/dto"
	"go-api/repository"
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
	user := s.userRepository.FindByClerkID(id)

	if user == nil {
		return nil
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Banned = banned

	return s.userRepository.Update(user)
}

func (s *UserService) DeleteUser(c fiber.Ctx, id string) error {
	user := s.userRepository.FindByClerkID(id)

	if user == nil {
		return nil
	}

	return s.userRepository.Delete(user)
}

func (s *UserService) GetUser(user *domain.User) (*dto.UserOutput, error) {
	output := dto.NewUserOutput(*user)
	return &output, nil
}
