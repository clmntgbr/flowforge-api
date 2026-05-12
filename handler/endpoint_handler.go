package handler

import (
	"flowforge-api/handler/context"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"flowforge-api/infrastructure/paginate"
	"flowforge-api/usecase/endpoint"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type EndpointHandler struct {
	listEndpointsUseCase  *endpoint.ListEndpointsUseCase
	createEndpointUseCase *endpoint.CreateEndpointUseCase
}

func NewEndpointHandler(
	listEndpointsUseCase *endpoint.ListEndpointsUseCase,
	createEndpointUseCase *endpoint.CreateEndpointUseCase,
) *EndpointHandler {
	return &EndpointHandler{
		listEndpointsUseCase:  listEndpointsUseCase,
		createEndpointUseCase: createEndpointUseCase,
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

	return c.JSON(paginate.NewPaginateResponse(endpoints, int(total), query))
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
		log.Println("error validating request: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	log.Println("request: ", request)

	_, err = h.createEndpointUseCase.Execute(c.Context(), activeOrganizationID, request)
	if err != nil {
		log.Println("error creating endpoint: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create endpoint",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *EndpointHandler) GetEndpointByID(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func (h *EndpointHandler) UpdateEndpoint(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
