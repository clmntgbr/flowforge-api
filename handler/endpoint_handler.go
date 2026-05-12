package handler

import (
	"flowforge-api/handler/context"
	"flowforge-api/infrastructure/paginate"
	"flowforge-api/usecase/endpoint"

	"github.com/gofiber/fiber/v3"
)

type EndpointHandler struct {
	listEndpointsUseCase *endpoint.ListEndpointsUseCase
}

func NewEndpointHandler(
	listEndpointsUseCase *endpoint.ListEndpointsUseCase,
) *EndpointHandler {
	return &EndpointHandler{
		listEndpointsUseCase: listEndpointsUseCase,
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
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
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
