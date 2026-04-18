package handler

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/service"

	"github.com/gofiber/fiber/v3"
)

type StepHandler struct {
	BaseHandler
	stepService *service.StepService
}

func NewStepHandler(stepService *service.StepService) *StepHandler {
	return &StepHandler{
		stepService: stepService,
	}
}

func (h *StepHandler) GetStepByID(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	stepUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidStepID)
	if err != nil {
		return err
	}

	step, err := h.stepService.GetStepByID(c, activeOrganizationID, stepUUID)
	if err != nil {
		return h.sendInternalError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(step)
}

func (h *StepHandler) UpdateStep(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.UpdateStepInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	stepUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidStepID)
	if err != nil {
		return err
	}

	_, err = h.stepService.UpdateStep(c, activeOrganizationID, stepUUID, req)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
