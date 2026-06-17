package handler

import (
	"flowforge-api/domain/entity"
	"flowforge-api/handler/context"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"flowforge-api/infrastructure/paginate"
	"flowforge-api/presenter"
	"flowforge-api/usecase/endpoint"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type EndpointHandler struct {
	listEndpointsUseCase     *endpoint.ListEndpointsUseCase
	createEndpointUseCase    *endpoint.CreateEndpointUseCase
	updateEndpointUseCase    *endpoint.UpdateEndpointUseCase
	getEndpointUseCase       *endpoint.GetEndpointUseCase
	importFromOpenAPIUseCase *endpoint.ImportFromOpenAPIUseCase
	endpointHasStepUseCase   *endpoint.EndpointHasStepUseCase
	deleteEndpointUseCase    *endpoint.DeleteEndpointUseCase
}

func NewEndpointHandler(
	listEndpointsUseCase *endpoint.ListEndpointsUseCase,
	createEndpointUseCase *endpoint.CreateEndpointUseCase,
	updateEndpointUseCase *endpoint.UpdateEndpointUseCase,
	getEndpointUseCase *endpoint.GetEndpointUseCase,
	importFromOpenAPIUseCase *endpoint.ImportFromOpenAPIUseCase,
	endpointHasStepUseCase *endpoint.EndpointHasStepUseCase,
	deleteEndpointUseCase *endpoint.DeleteEndpointUseCase,
) *EndpointHandler {
	return &EndpointHandler{
		listEndpointsUseCase:     listEndpointsUseCase,
		createEndpointUseCase:    createEndpointUseCase,
		updateEndpointUseCase:    updateEndpointUseCase,
		getEndpointUseCase:       getEndpointUseCase,
		importFromOpenAPIUseCase: importFromOpenAPIUseCase,
		endpointHasStepUseCase:   endpointHasStepUseCase,
		deleteEndpointUseCase:    deleteEndpointUseCase,
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

	_, err = h.createEndpointUseCase.Execute(c.Context(), activeOrganizationID, request, []entity.Tag{})
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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
		})
	}

	endpoint, err := h.getEndpointUseCase.Execute(c.Context(), activeOrganizationID, endpointUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get endpoint",
		})
	}

	return c.JSON(presenter.NewEndpointDetailResponse(endpoint))
}

func (h *EndpointHandler) UpdateEndpoint(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
		})
	}

	var request endpointDTO.UpdateEndpointInput
	if err := c.Bind().JSON(&request); err != nil {
		fmt.Println("error1", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		fmt.Println("error2", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	_, err = h.updateEndpointUseCase.Execute(c.Context(), activeOrganizationID, endpointUUID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update endpoint",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *EndpointHandler) ImportEndpoints(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var request endpointDTO.ImportEndpointsInput
	if err := c.Bind().Form(&request); err != nil {
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

	err = h.importFromOpenAPIUseCase.Execute(c.Context(), activeOrganizationID, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to import endpoints",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Endpoints imported successfully",
	})
}

func (h *EndpointHandler) DeleteEndpoint(c fiber.Ctx) error {
	_, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
		})
	}

	hasSteps, err := h.endpointHasStepUseCase.Execute(c.Context(), endpointUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check if endpoint has steps",
		})
	}

	if hasSteps {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Endpoint has steps",
		})
	}

	err = h.deleteEndpointUseCase.Execute(c.Context(), endpointUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete endpoint",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
