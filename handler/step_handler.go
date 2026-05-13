package handler

import (
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
}

func NewStepHandler(getStepUseCase *step.GetStepUseCase, updateStepUseCase *step.UpdateStepUseCase) *StepHandler {
	return &StepHandler{getStepUseCase: getStepUseCase, updateStepUseCase: updateStepUseCase}
}

func (h *StepHandler) GetStepByID(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	stepID := c.Params("id")
	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid step ID",
		})
	}

	workflowID := c.Params("workflowId")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	step, err := h.getStepUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID, stepUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get step",
		})
	}

	return c.JSON(presenter.NewStepDetailResponse(step))
}

func (h *StepHandler) UpdateStep(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	stepID := c.Params("id")
	stepUUID, err := uuid.Parse(stepID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid step ID",
		})
	}

	workflowID := c.Params("workflowId")
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workflow ID",
		})
	}

	var request stepDTO.UpdateStepInput
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

	_, err = h.updateStepUseCase.Execute(c.Context(), activeOrganizationID, workflowUUID, stepUUID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update step",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *StepHandler) DeleteStep(c fiber.Ctx) error {
	return nil
}
