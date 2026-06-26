package handler

import (
	"flowforge-api/handler/context"
	variableDTO "flowforge-api/infrastructure/variable"
	"flowforge-api/presenter"
	"flowforge-api/usecase/variable"
	"flowforge-api/usecase/workflow"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type VariableHandler struct {
	getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase
	createVariableUseCase           *variable.CreateVariableUseCase
	getWorkflowUseCase              *workflow.GetWorkflowUseCase
	searchVariablesPathUseCase      *variable.SearchVariablesPathUseCase
}

func NewVariableHandler(
	getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase,
	createVariableUseCase *variable.CreateVariableUseCase,
	getWorkflowUseCase *workflow.GetWorkflowUseCase,
	searchVariablesPathUseCase *variable.SearchVariablesPathUseCase,
) *VariableHandler {
	return &VariableHandler{
		getVariablesByWorkflowIDUseCase: getVariablesByWorkflowIDUseCase,
		createVariableUseCase:           createVariableUseCase,
		getWorkflowUseCase:              getWorkflowUseCase,
		searchVariablesPathUseCase:      searchVariablesPathUseCase,
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

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	variables, err := h.getVariablesByWorkflowIDUseCase.Execute(c.Context(), workflow.ID)
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

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
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

	_, err = h.createVariableUseCase.Execute(c.Context(), workflow.ID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create variable",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *VariableHandler) SearchVariablesPath(c fiber.Ctx) error {
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

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	var request variableDTO.SearchVariablesPathInput
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

	variables, err := h.searchVariablesPathUseCase.Execute(c.Context(), workflow.ID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"paths": variables,
	})
}
