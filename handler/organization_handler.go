package handler

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/service"

	"github.com/gofiber/fiber/v3"
)

type OrganizationHandler struct {
	BaseHandler
	organizationService service.OrganizationServiceInterface
}

func NewOrganizationHandler(organizationService service.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: organizationService,
	}
}

func (h *OrganizationHandler) GetOrganizations(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	output, err := h.organizationService.GetOrganizations(c, user, activeOrganizationID)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

func (h *OrganizationHandler) GetOrganizationByID(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	organizationUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidOrganizationID)
	if err != nil {
		return err
	}

	organization, err := h.organizationService.GetOrganizationByID(c, user, organizationUUID)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(organization)
}

func (h *OrganizationHandler) CreateOrganization(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.CreateOrganizationInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	_, err = h.organizationService.CreateOrganization(c, user, req.Name)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *OrganizationHandler) UpdateOrganization(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.UpdateOrganizationInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	organizationUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidOrganizationID)
	if err != nil {
		return err
	}

	_, err = h.organizationService.UpdateOrganization(c, user, organizationUUID, req)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func (h *OrganizationHandler) ActivateOrganization(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	organizationUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidOrganizationID)
	if err != nil {
		return err
	}

	_, err = h.organizationService.ActivateOrganization(c.Context(), user.ID, organizationUUID)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
