package service

import (
	"go-api/dto"
	"go-api/repository"

	"github.com/gofiber/fiber/v3"
)

type WebhookClerkService struct {
	userRepository *repository.UserRepository
	userService    *UserService
}

func NewWebhookClerkService(userRepository *repository.UserRepository, userService *UserService) *WebhookClerkService {
	return &WebhookClerkService{
		userRepository: userRepository,
		userService:    userService,
	}
}

func (s *WebhookClerkService) CreateUser(c fiber.Ctx, data dto.ClerkUserCreated) error {
	user := s.userRepository.FindByClerkID(data.ID)

	if user != nil {
		return nil
	}

	_, err := s.userService.CreateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned)
	if err != nil {
		return err
	}
	return nil
}

func (s *WebhookClerkService) UpdateUser(c fiber.Ctx, data dto.ClerkUserUpdated) error {
	user := s.userRepository.FindByClerkID(data.ID)

	if user == nil {
		return nil
	}

	return s.userService.UpdateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned)
}

func (s *WebhookClerkService) DeleteUser(c fiber.Ctx, data dto.ClerkUserDeleted) error {
	return s.userService.DeleteUser(c, data.ID)
}
