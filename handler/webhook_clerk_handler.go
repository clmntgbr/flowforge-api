package handler

import (
	"encoding/json"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/usecase"
	"log"

	"github.com/gofiber/fiber/v3"
)

type WebhookClerkHandler struct {
	BaseHandler
	createUserUsecase *usecase.CreateUserUsecase
	updateUserUsecase *usecase.UpdateUserUsecase
	deleteUserUsecase *usecase.DeleteUserUsecase
}

func NewWebhookClerkHandler(
	createUserUsecase *usecase.CreateUserUsecase,
	updateUserUsecase *usecase.UpdateUserUsecase,
	deleteUserUsecase *usecase.DeleteUserUsecase,
) *WebhookClerkHandler {
	return &WebhookClerkHandler{
		createUserUsecase: createUserUsecase,
		updateUserUsecase: updateUserUsecase,
		deleteUserUsecase: deleteUserUsecase,
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

		if _, _, err := h.createUserUsecase.CreateUser(c.Context(), data.ID, data.FirstName, data.LastName, *data.Banned); err != nil {
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

		if err := h.updateUserUsecase.UpdateUser(c.Context(), data.ID, data.FirstName, data.LastName, *data.Banned); err != nil {
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

		if err := h.deleteUserUsecase.DeleteUser(c.Context(), data.ID); err != nil {
			return h.sendInternalError(c, err)
		}

		return c.SendStatus(fiber.StatusNoContent)

	default:
		log.Printf("Unhandled event type: %s", clerkEvent.Type)
		return c.SendStatus(fiber.StatusOK)
	}
}
