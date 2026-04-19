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
	connexionService service.ConnexionServiceInterface
	workflowService  service.WorkflowServiceInterface
}

func NewConnexionHandler(connexionService service.ConnexionServiceInterface, workflowService service.WorkflowServiceInterface) *ConnexionHandler {
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

	_, err = h.workflowService.GetWorkflowByID(c, activeOrganizationID, req.WorkflowID)
	if err != nil {
		return h.sendError(c, err)
	}

	connexion, err := h.connexionService.CreateConnexion(c, req.WorkflowID, req)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(connexion)
}

func (h *ConnexionHandler) DeleteConnexion(c fiber.Ctx) error {
	connexionUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidConnexionID)
	if err != nil {
		return err
	}

	_, err = h.connexionService.DeleteConnexion(c.Context(), connexionUUID)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
