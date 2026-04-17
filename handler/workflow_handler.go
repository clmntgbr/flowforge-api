package handler

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/service"

	"github.com/gofiber/fiber/v3"
)

type WorkflowHandler struct {
	BaseHandler
	workflowService *service.WorkflowService
}

func NewWorkflowHandler(workflowService *service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
	}
}

func (h *WorkflowHandler) GetWorkflows(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	query, err := h.bindPaginateQuery(c)
	if err != nil {
		return err
	}

	output, err := h.workflowService.GetWorkflows(c, activeOrganizationID, query)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

func (h *WorkflowHandler) CreateWorkflow(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.CreateWorkflowInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	_, err = h.workflowService.CreateWorkflow(c, activeOrganizationID, req)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *WorkflowHandler) GetWorkflowByID(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	workflowUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidWorkflowID)
	if err != nil {
		return err
	}

	workflow, err := h.workflowService.GetWorkflowByID(c, activeOrganizationID, workflowUUID)
	if err != nil {
		return h.sendInternalError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(workflow)
}

func (h *WorkflowHandler) UpdateWorkflow(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.UpdateWorkflowInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	workflowUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidWorkflowID)
	if err != nil {
		return err
	}

	_, err = h.workflowService.UpdateWorkflow(c, activeOrganizationID, workflowUUID, req)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
