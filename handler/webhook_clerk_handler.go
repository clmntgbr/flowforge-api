package handler

import (
	"encoding/json"
	"fmt"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/repository"
	"forgeflow-api/service"
	"log"

	"github.com/gofiber/fiber/v3"
)

type WebhookClerkHandler struct {
	BaseHandler
	userService    *service.UserService
	projectService *service.ProjectService
	userRepository *repository.UserRepository
}

func NewWebhookClerkHandler(userService *service.UserService, projectService *service.ProjectService, userRepository *repository.UserRepository) *WebhookClerkHandler {
	return &WebhookClerkHandler{
		userService:    userService,
		projectService: projectService,
		userRepository: userRepository,
	}
}

func (h *WebhookClerkHandler) Handle(c fiber.Ctx) error {
	clerkEvent := c.Locals("payload").(dto.ClerkEvent)

	switch clerkEvent.Type {
	case "user.created":
		var data dto.ClerkUserCreated
		if err := json.Unmarshal(clerkEvent.Data, &data); err != nil {
			return h.sendBadRequest(c, errors.ErrInvalidRequestBody)
		}

		if err := h.validate(c, &data); err != nil {
			return err
		}

		if err := h.CreateUser(c, data); err != nil {
			return h.sendInternalError(c, err)
		}

		return c.SendStatus(fiber.StatusCreated)

	case "user.updated":
		var data dto.ClerkUserUpdated
		if err := json.Unmarshal(clerkEvent.Data, &data); err != nil {
			return h.sendBadRequest(c, errors.ErrInvalidRequestBody)
		}

		if err := h.validate(c, &data); err != nil {
			return err
		}

		if err := h.UpdateUser(c, data); err != nil {
			return h.sendInternalError(c, err)
		}

		return c.SendStatus(fiber.StatusNoContent)

	case "user.deleted":
		var data dto.ClerkUserDeleted
		if err := json.Unmarshal(clerkEvent.Data, &data); err != nil {
			return h.sendBadRequest(c, errors.ErrInvalidRequestBody)
		}

		if err := h.validate(c, &data); err != nil {
			return err
		}

		if err := h.DeleteUser(c, data); err != nil {
			return h.sendInternalError(c, err)
		}

		return c.SendStatus(fiber.StatusNoContent)

	default:
		log.Printf("Unhandled event type: %s", clerkEvent.Type)
		return c.SendStatus(fiber.StatusOK)
	}
}

func (h *WebhookClerkHandler) CreateUser(c fiber.Ctx, data dto.ClerkUserCreated) error {
	user, err := h.userRepository.FindByClerkID(data.ID)
	if err != nil {
		fmt.Println("Error finding user", err)
		return err
	}

	if user != nil {
		fmt.Println("User already exists")
		return nil
	}

	user, err = h.userService.CreateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned)
	if err != nil {
		fmt.Println("Error creating user", err)
		return err
	}

	project, err := h.projectService.CreateProject(c, user, "Default Project")
	if err != nil {
		fmt.Println("Error creating project", err)
		return err
	}

	user.ActiveProjectID = &project.ID
	if err := h.userRepository.Update(user); err != nil {
		fmt.Println("Error updating user", err)
		return err
	}

	return nil
}

func (h *WebhookClerkHandler) UpdateUser(c fiber.Ctx, data dto.ClerkUserUpdated) error {
	user, err := h.userRepository.FindByClerkID(data.ID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	return h.userService.UpdateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned)
}

func (h *WebhookClerkHandler) DeleteUser(c fiber.Ctx, data dto.ClerkUserDeleted) error {
	return h.userService.DeleteUser(c, data.ID)
}
