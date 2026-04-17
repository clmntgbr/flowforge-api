package handler

import (
	"forgeflow-api/ctxutil"
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
