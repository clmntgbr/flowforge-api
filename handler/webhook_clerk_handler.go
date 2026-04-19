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
	userService         service.UserServiceInterface
	organizationService service.OrganizationServiceInterface
	userRepository      repository.UserRepositoryInterface
}

func NewWebhookClerkHandler(userService service.UserServiceInterface, organizationService service.OrganizationServiceInterface, userRepository repository.UserRepositoryInterface) *WebhookClerkHandler {
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
		log.Printf("Error finding user by Clerk ID %s: %v", data.ID, err)
		return errors.ErrUserFailedToCreate
	}

	if user != nil {
		log.Printf("User with Clerk ID %s already exists, skipping creation", data.ID)
		return nil
	}

	user, err = h.userService.CreateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned)
	if err != nil {
		log.Printf("Error creating user with Clerk ID %s: %v", data.ID, err)
		return errors.ErrUserFailedToCreate
	}

	organization, err := h.organizationService.CreateOrganization(c, user, "Default Organization")
	if err != nil {
		log.Printf("Error creating default organization for user %s: %v", user.ID, err)
		return errors.ErrOrganizationFailedToCreate
	}

	organizationID, err := uuid.Parse(organization.ID)
	if err != nil {
		log.Printf("Error parsing organization UUID %s: %v", organization.ID, err)
		return errors.ErrInvalidOrganizationUUID
	}

	user.ActiveOrganizationID = &organizationID
	if err := h.userRepository.Update(user); err != nil {
		log.Printf("Error updating user %s with active organization: %v", user.ID, err)
		return errors.ErrUserFailedToCreate
	}

	log.Printf("Successfully created user with Clerk ID %s and organization %s", data.ID, organizationID)
	return nil
}

func (h *WebhookClerkHandler) UpdateUser(c fiber.Ctx, data dto.ClerkUserUpdated) error {
	user, err := h.userRepository.FindByClerkID(data.ID)
	if err != nil {
		log.Printf("Error finding user by Clerk ID %s: %v", data.ID, err)
		return err
	}

	if user == nil {
		log.Printf("User with Clerk ID %s not found for update", data.ID)
		return errors.ErrUserNotFound
	}

	if err := h.userService.UpdateUser(c, data.ID, data.FirstName, data.LastName, *data.Banned); err != nil {
		log.Printf("Error updating user with Clerk ID %s: %v", data.ID, err)
		return err
	}

	log.Printf("Successfully updated user with Clerk ID %s", data.ID)
	return nil
}

func (h *WebhookClerkHandler) DeleteUser(c fiber.Ctx, data dto.ClerkUserDeleted) error {
	if err := h.userService.DeleteUser(c, data.ID); err != nil {
		log.Printf("Error deleting user with Clerk ID %s: %v", data.ID, err)
		return err
	}

	log.Printf("Successfully deleted user with Clerk ID %s", data.ID)
	return nil
}
