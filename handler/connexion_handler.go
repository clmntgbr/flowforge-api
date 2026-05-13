package handler

import (
	"flowforge-api/handler/context"
	connexionDTO "flowforge-api/infrastructure/connexion"
	"flowforge-api/usecase/connexion"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ConnexionHandler struct {
	createConnexionUseCase *connexion.CreateConnexionUseCase
	deleteConnexionUseCase *connexion.DeleteConnexionUseCase
}

func NewConnexionHandler(
	createConnexionUseCase *connexion.CreateConnexionUseCase,
	deleteConnexionUseCase *connexion.DeleteConnexionUseCase,
) *ConnexionHandler {
	return &ConnexionHandler{
		createConnexionUseCase: createConnexionUseCase,
		deleteConnexionUseCase: deleteConnexionUseCase,
	}
}

func (h *ConnexionHandler) CreateConnexion(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var request connexionDTO.CreateConnexionInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	_, err = h.createConnexionUseCase.Execute(c.Context(), activeOrganizationID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create connexion",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *ConnexionHandler) DeleteConnexion(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
		})
	}

	err = h.deleteConnexionUseCase.Execute(c.Context(), activeOrganizationID, endpointUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete connexion",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
