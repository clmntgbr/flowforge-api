package handler

import (
	"flowforge-api/handler/context"
	"flowforge-api/presenter"
	"flowforge-api/usecase/variable"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type VariableHandler struct {
	getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase
}

func NewVariableHandler(getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase) *VariableHandler {
	return &VariableHandler{getVariablesByWorkflowIDUseCase: getVariablesByWorkflowIDUseCase}
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
