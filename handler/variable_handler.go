package handler

import (
	"errors"
	"flowforge-api/handler/context"
	variableDTO "flowforge-api/infrastructure/variable"
	"flowforge-api/presenter"
	"flowforge-api/usecase/variable"
	"flowforge-api/usecase/workflow"
	"strings"

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
	deleteVariableUseCase           *variable.DeleteVariableUseCase
}

func NewVariableHandler(
	getVariablesByWorkflowIDUseCase *variable.GetVariablesByWorkflowIDUseCase,
	createVariableUseCase *variable.CreateVariableUseCase,
	getWorkflowUseCase *workflow.GetWorkflowUseCase,
	searchVariablesPathUseCase *variable.SearchVariablesPathUseCase,
	getVariableByIDUseCase *variable.GetVariableByIDUseCase,
	updateVariableUseCase *variable.UpdateVariableUseCase,
	deleteVariableUseCase *variable.DeleteVariableUseCase,
) *VariableHandler {
	return &VariableHandler{
		getVariablesByWorkflowIDUseCase: getVariablesByWorkflowIDUseCase,
		createVariableUseCase:           createVariableUseCase,
		getWorkflowUseCase:              getWorkflowUseCase,
		searchVariablesPathUseCase:      searchVariablesPathUseCase,
		getVariableByIDUseCase:          getVariableByIDUseCase,
		updateVariableUseCase:           updateVariableUseCase,
		deleteVariableUseCase:           deleteVariableUseCase,
	}
}

func (h *VariableHandler) GetVariablesByWorkflowID(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	variables, err := h.getVariablesByWorkflowIDUseCase.Execute(c.Context(), workflow.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	return c.JSON(presenter.NewVariableListResponses(variables))
}

func (h *VariableHandler) CreateVariable(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	var request variableDTO.CreateVariableInput
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

	_, err = h.createVariableUseCase.Execute(c.Context(), workflow.ID, request)
	if err != nil {
		if strings.Contains(err.Error(), "variable key already exists in workflow") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Variable key already exists in this workflow",
				"errors":  err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create variable",
			"errors":  err.Error(),
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
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	var request variableDTO.SearchVariablesPathInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	if err := c.Bind().Query(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query parameters",
			"errors":  err.Error(),
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
			"errors":  err.Error(),
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
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	variableID := c.Params("variableId")
	variableUUID, err := uuid.Parse(variableID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid variable ID",
			"errors":  err.Error(),
		})
	}

	variable, err := h.getVariableByIDUseCase.Execute(c.Context(), workflow.ID, variableUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	return c.JSON(presenter.NewVariableDetailResponses(variable))
}

func (h *VariableHandler) UpdateVariable(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	variableID := c.Params("variableId")
	variableUUID, err := uuid.Parse(variableID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid variable ID",
			"errors":  err.Error(),
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
		if strings.Contains(err.Error(), "variable key already exists in workflow") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Variable key already exists in this workflow",
				"errors":  err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update variable",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func (h *VariableHandler) DeleteVariable(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("id")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	workflow, err := h.getWorkflowUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	variableID := c.Params("variableId")
	variableUUID, err := uuid.Parse(variableID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid variable ID",
			"errors":  err.Error(),
		})
	}

	err = h.deleteVariableUseCase.Execute(c.Context(), workflow.ID, variableUUID)
	if err != nil {
		var inUseErr *variable.VariableInUseError
		if errors.As(err, &inUseErr) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Variable is used in workflow steps",
				"errors":  err.Error(),
				"steps":   presenter.NewStepDetailResponses(inUseErr.Steps),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete variable",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
