package handler

import (
	"go-api/errors"
	"go-api/validator"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BaseHandler struct{}

func (h *BaseHandler) bindAndValidate(c fiber.Ctx, req interface{}) error {
	if err := c.Bind().JSON(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.ErrInvalidRequestBody.Error(),
		})
		return errors.ErrValidationFailed
	}

	if err := validator.ValidateStruct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.ErrInvalidRequestBody.Error(),
			"errors":  validator.FormatValidationErrors(err),
		})

		return errors.ErrValidationFailed
	}

	return nil
}

func (h *BaseHandler) validate(c fiber.Ctx, req interface{}) error {
	if err := validator.ValidateStruct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors.ErrInvalidRequestBody.Error(),
			"errors":  validator.FormatValidationErrors(err),
		})

		return errors.ErrValidationFailed
	}
	return nil
}

func (h *BaseHandler) sendUnauthorized(c fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": errors.ErrUserNotAuthenticated.Error(),
	})
}

func (h *BaseHandler) sendInternalError(c fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func (h *BaseHandler) sendBadRequest(c fiber.Ctx, message error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": message.Error(),
	})
}

func (h *BaseHandler) sendNotFound(c fiber.Ctx, err error) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func (h *BaseHandler) parseUUIDParam(c fiber.Ctx, param string, customErr error) (uuid.UUID, error) {
	raw := c.Params(param)
	if raw == "" {
		h.sendBadRequest(c, customErr)
		return uuid.Nil, customErr
	}
	parsed, err := uuid.Parse(raw)
	if err != nil {
		h.sendBadRequest(c, customErr)
		return uuid.Nil, customErr
	}
	return parsed, nil
}
