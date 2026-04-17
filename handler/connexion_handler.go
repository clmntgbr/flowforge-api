package handler

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/service"

	"github.com/gofiber/fiber/v3"
)

type ConnexionHandler struct {
	BaseHandler
	connexionService *service.ConnexionService
	workflowService  *service.WorkflowService
}

func NewConnexionHandler(connexionService *service.ConnexionService, workflowService *service.WorkflowService) *ConnexionHandler {
	return &ConnexionHandler{
		connexionService: connexionService,
		workflowService:  workflowService,
	}
}

func (h *ConnexionHandler) CreateConnexion(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.CreateConnexionInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	workflowUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidWorkflowID)
	if err != nil {
		return err
	}

	_, err = h.workflowService.GetWorkflowByID(c, activeOrganizationID, workflowUUID)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	_, err = h.connexionService.CreateConnexion(c, workflowUUID, req)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *ConnexionHandler) DeleteConnexion(c fiber.Ctx) error {
	connexionUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidConnexionID)
	if err != nil {
		return err
	}

	_, err = h.connexionService.DeleteConnexion(c.Context(), connexionUUID)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
