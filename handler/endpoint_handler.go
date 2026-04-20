package handler

import (
	"fmt"
	"forgeflow-api/ctxutil"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
	"forgeflow-api/service"

	"github.com/gofiber/fiber/v3"
)

type EndpointHandler struct {
	BaseHandler
	endpointService service.EndpointServiceInterface
}

func NewEndpointHandler(endpointService service.EndpointServiceInterface) *EndpointHandler {
	return &EndpointHandler{
		endpointService: endpointService,
	}
}

func (h *EndpointHandler) GetEndpoints(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	query, err := h.bindPaginateQuery(c)
	if err != nil {
		return h.sendBadRequest(c, err)
	}

	output, err := h.endpointService.GetEndpoints(c, activeOrganizationID, query)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

func (h *EndpointHandler) CreateEndpoint(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.CreateEndpointInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		fmt.Println("Error binding and validating request", err)
		fmt.Println("Response", response)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	_, err = h.endpointService.CreateEndpoint(c, activeOrganizationID, req)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *EndpointHandler) GetEndpointByID(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	endpointUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidEndpointID)
	if err != nil {
		return err
	}

	endpoint, err := h.endpointService.GetEndpointByID(c, activeOrganizationID, endpointUUID)
	if err != nil {
		return h.sendError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(endpoint)
}

func (h *EndpointHandler) UpdateEndpoint(c fiber.Ctx) error {
	activeOrganizationID, err := ctxutil.GetOrganizationID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.UpdateEndpointInput
	err, response := h.bindAndValidate(c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	endpointUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidEndpointID)
	if err != nil {
		return err
	}

	_, err = h.endpointService.UpdateEndpoint(c, activeOrganizationID, endpointUUID, req)
	if err != nil {
		return h.sendError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
