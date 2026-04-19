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
	"github.com/google/uuid"
)

type WebhookClerkHandler struct {
	BaseHandler
	userService         *service.UserService
	organizationService *service.OrganizationService
	userRepository      *repository.UserRepository
}

func NewWebhookClerkHandler(userService *service.UserService, organizationService *service.OrganizationService, userRepository *repository.UserRepository) *WebhookClerkHandler {
	return &WebhookClerkHandler{
		userService:         userService,
		organizationService: organizationService,
		userRepository:      userRepository,
	}
}

func (h *WebhookClerkHandler) Handle(c fiber.Ctx) error {
	clerkEvent := c.Locals("payload").(dto.ClerkEvent)

	fmt.Println("Clerk event type", clerkEvent.Type)

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
			return h.sendError(c, err)
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
			return h.sendError(c, err)
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
			return h.sendError(c, err)
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
		return err
	}

	if user != nil {
		return nil
	}

	user, err = h.userService.CreateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned)
	if err != nil {
		return err
	}

	organization, err := h.organizationService.CreateOrganization(c, user, "Default Organization")
	if err != nil {
		return err
	}

	organizationID, err := uuid.Parse(organization.ID)
	if err != nil {
		return errors.ErrOrganizationFailedToCreate
	}

	user.ActiveOrganizationID = &organizationID
	if err := h.userRepository.Update(user); err != nil {
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
