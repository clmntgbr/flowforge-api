package handler

import (
	"flowforge-api/handler/context"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"flowforge-api/infrastructure/paginate"
	"flowforge-api/presenter"
	"flowforge-api/usecase/endpoint"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type EndpointHandler struct {
	listEndpointsUseCase  *endpoint.ListEndpointsUseCase
	createEndpointUseCase *endpoint.CreateEndpointUseCase
	updateEndpointUseCase *endpoint.UpdateEndpointUseCase
	getEndpointUseCase    *endpoint.GetEndpointUseCase
}

func NewEndpointHandler(
	listEndpointsUseCase *endpoint.ListEndpointsUseCase,
	createEndpointUseCase *endpoint.CreateEndpointUseCase,
	updateEndpointUseCase *endpoint.UpdateEndpointUseCase,
	getEndpointUseCase *endpoint.GetEndpointUseCase,
) *EndpointHandler {
	return &EndpointHandler{
		listEndpointsUseCase:  listEndpointsUseCase,
		createEndpointUseCase: createEndpointUseCase,
		updateEndpointUseCase: updateEndpointUseCase,
		getEndpointUseCase:    getEndpointUseCase,
	}
}

func (h *EndpointHandler) GetEndpoints(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var query paginate.PaginateQuery
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	query.Normalize()

	endpoints, total, err := h.listEndpointsUseCase.Execute(c.Context(), activeOrganizationID, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(paginate.NewPaginateResponse(presenter.NewEndpointListResponses(endpoints), int(total), query))
}

func (h *EndpointHandler) CreateEndpoint(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var request endpointDTO.CreateEndpointInput
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

	_, err = h.createEndpointUseCase.Execute(c.Context(), activeOrganizationID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create endpoint",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *EndpointHandler) GetEndpointByID(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		log.Println("error getting organization ID: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		log.Println("error parsing endpoint ID: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
		})
	}

	endpoint, err := h.getEndpointUseCase.Execute(c.Context(), activeOrganizationID, endpointUUID)
	if err != nil {
		log.Println("error getting endpoint: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get endpoint",
		})
	}

	return c.JSON(presenter.NewEndpointDetailResponse(endpoint))
}

func (h *EndpointHandler) UpdateEndpoint(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		log.Println("error getting organization ID: ", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		log.Println("error parsing endpoint ID: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
		})
	}

	var request endpointDTO.UpdateEndpointInput
	if err := c.Bind().JSON(&request); err != nil {
		log.Println("error binding request body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		log.Println("error validating request body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	_, err = h.updateEndpointUseCase.Execute(c.Context(), activeOrganizationID, endpointUUID, request)
	if err != nil {
		log.Println("error updating endpoint: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update endpoint",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
