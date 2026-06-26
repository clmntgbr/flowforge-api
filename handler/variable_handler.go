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
	getVariableByIDUseCase          *variable.GetVariableByIDUseCase
	updateVariableUseCase           *variable.UpdateVariableUseCase
}

func NewVariableHandler(
	getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase,
	createVariableUseCase *variable.CreateVariableUseCase,
	getWorkflowUseCase *workflow.GetWorkflowUseCase,
	searchVariablesPathUseCase *variable.SearchVariablesPathUseCase,
	getVariableByIDUseCase *variable.GetVariableByIDUseCase,
	updateVariableUseCase *variable.UpdateVariableUseCase,
) *VariableHandler {
	return &VariableHandler{
		getVariablesByWorkflowIDUseCase: getVariablesByWorkflowIDUseCase,
		createVariableUseCase:           createVariableUseCase,
		getWorkflowUseCase:              getWorkflowUseCase,
		searchVariablesPathUseCase:      searchVariablesPathUseCase,
		getVariableByIDUseCase:          getVariableByIDUseCase,
		updateVariableUseCase:           updateVariableUseCase,
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

	if err := c.Bind().Query(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query parameters",
		})
	}

	request.Normalize()

	if err := validator.New().Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	paths, total, err := h.searchVariablesPathUseCase.Execute(c.Context(), workflow.ID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	totalPages := (total + request.Limit - 1) / request.Limit

	return c.JSON(fiber.Map{
		"paths":      paths,
		"total":      total,
		"page":       request.Page,
		"limit":      request.Limit,
		"totalPages": totalPages,
	})
}

func (h *VariableHandler) GetVariableByID(c fiber.Ctx) error {
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

	variableID := c.Params("variableId")
	variableUUID, err := uuid.Parse(variableID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid variable ID",
		})
	}

	variable, err := h.getVariableByIDUseCase.Execute(c.Context(), workflow.ID, variableUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(presenter.NewVariableResponse(variable))
}

func (h *VariableHandler) UpdateVariable(c fiber.Ctx) error {
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

	variableID := c.Params("variableId")
	variableUUID, err := uuid.Parse(variableID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid variable ID",
		})
	}

	var request variableDTO.UpdateVariableInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	if err := validator.New().Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	_, err = h.updateVariableUseCase.Execute(c.Context(), workflow.ID, variableUUID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update variable",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
