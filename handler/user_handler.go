package handler

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/usecase"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	BaseHandler
	getUserUsecase *usecase.GetUserUsecase
}

func NewUserHandler(getUserUsecase *usecase.GetUserUsecase) *UserHandler {
	return &UserHandler{
		getUserUsecase: getUserUsecase,
	}
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	output, err := h.getUserUsecase.GetUser(user)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
