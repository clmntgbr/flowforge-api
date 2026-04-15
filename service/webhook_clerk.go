package service

import (
	"forgeflow-api/dto"
	"forgeflow-api/repository"

	"github.com/gofiber/fiber/v3"
)

type WebhookClerkService struct {
	userRepository *repository.UserRepository
	userService    *UserService
	projectService *ProjectService
}

func NewWebhookClerkService(userRepository *repository.UserRepository, userService *UserService, projectService *ProjectService) *WebhookClerkService {
	return &WebhookClerkService{
		userRepository: userRepository,
		userService:    userService,
		projectService: projectService,
	}
}

func (s *WebhookClerkService) CreateUser(c fiber.Ctx, data dto.ClerkUserCreated) error {
	user := s.userRepository.FindByClerkID(data.ID)

	if user != nil {
		return nil
	}

	user, err := s.userService.CreateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned)
	if err != nil {
		return err
	}

	_, err = s.projectService.CreateProject(c, user, "Default Project")
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
