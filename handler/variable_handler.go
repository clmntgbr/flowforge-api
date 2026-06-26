package handler

import (
	"flowforge-api/handler/context"
	variableDTO "flowforge-api/infrastructure/variable"
	"flowforge-api/presenter"
	"flowforge-api/usecase/variable"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type VariableHandler struct {
	getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase
	createVariableUseCase           *variable.CreateVariableUseCase
}

func NewVariableHandler(
	getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase,
	createVariableUseCase *variable.CreateVariableUseCase,
) *VariableHandler {
	return &VariableHandler{
		getVariablesByWorkflowIDUseCase: getVariablesByWorkflowIDUseCase,
		createVariableUseCase:           createVariableUseCase,
	}
}

func (h *VariableHandler) GetVariablesByWorkflowID(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	variables, err := h.getVariablesByWorkflowIDUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(presenter.NewVariableResponses(variables))
}

func (h *VariableHandler) CreateVariable(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	var request variableDTO.CreateVariableInput
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

	_, err = h.createVariableUseCase.Execute(c.Context(), workflowUUID, activeOrganizationID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create variable",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
