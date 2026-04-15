package handler

import (
	"encoding/json"
	"go-api/dto"
	"go-api/errors"
	"go-api/service"
	"log"

	"github.com/gofiber/fiber/v3"
)

type WebhookClerkHandler struct {
	BaseHandler
	webhookClerkService *service.WebhookClerkService
}

func NewWebhookClerkHandler(webhookClerkService *service.WebhookClerkService) *WebhookClerkHandler {
	return &WebhookClerkHandler{
		webhookClerkService: webhookClerkService,
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

		if err := h.webhookClerkService.CreateUser(c, data); err != nil {
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

		if err := h.webhookClerkService.UpdateUser(c, data); err != nil {
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

		if err := h.webhookClerkService.DeleteUser(c, data); err != nil {
			return h.sendInternalError(c, err)
		}

		return c.SendStatus(fiber.StatusNoContent)

	default:
		log.Printf("Unhandled event type: %s", clerkEvent.Type)
		return c.SendStatus(fiber.StatusOK)
	}
}
