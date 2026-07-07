package handler

import (
	"errors"
	"flowforge-api/handler/context"
	stepDTO "flowforge-api/infrastructure/step"
	"flowforge-api/presenter"
	"flowforge-api/usecase/step"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type StepHandler struct {
	getStepUseCase    *step.GetStepUseCase
	updateStepUseCase *step.UpdateStepUseCase
	deleteStepUseCase *step.DeleteStepUseCase
}

func NewStepHandler(getStepUseCase *step.GetStepUseCase, updateStepUseCase *step.UpdateStepUseCase, deleteStepUseCase *step.DeleteStepUseCase) *StepHandler {
	return &StepHandler{getStepUseCase: getStepUseCase, updateStepUseCase: updateStepUseCase, deleteStepUseCase: deleteStepUseCase}
}

func (h *StepHandler) GetStepByID(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	stepID := c.Params("id")
	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid step ID",
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("workflowId")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	step, err := h.getStepUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID, stepUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get step",
			"errors":  err.Error(),
		})
	}

	return c.JSON(presenter.NewStepDetailResponse(step))
}

func (h *StepHandler) UpdateStep(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	stepID := c.Params("id")
	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid step ID",
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("workflowId")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	var request stepDTO.UpdateStepInput
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

	_, err = h.updateStepUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID, stepUUID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update step",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *StepHandler) DeleteStep(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	stepID := c.Params("id")
	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid step ID",
			"errors":  err.Error(),
		})
	}

	workflowID := c.Params("workflowId")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
			"errors":  err.Error(),
		})
	}

	err = h.deleteStepUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID, stepUUID)
	if err != nil {
		var inUseErr *step.StepInUseError
		if errors.As(err, &inUseErr) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message":   "Step is used in workflow variables",
				"errors":    err.Error(),
				"variables": presenter.NewVariableDetailResponsesList(inUseErr.Variables),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete step",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Step deleted successfully",
	})
}
