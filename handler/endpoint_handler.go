package handler

import (
	"errors"
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
			"errors":  err.Error(),
		})
	}

	var query endpointDTO.PaginateEndpointQuery
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}
	query.Normalize()

	endpoints, total, err := h.listEndpointsUseCase.Execute(c.Context(), activeOrganizationID, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	return c.JSON(paginate.NewPaginateResponse(presenter.NewEndpointListResponses(endpoints), int(total), query.PaginateQuery))
}

func (h *EndpointHandler) CreateEndpoint(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	var request endpointDTO.CreateEndpointInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
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
			"errors":  err.Error(),
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
			"errors":  err.Error(),
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
			"errors":  err.Error(),
		})
	}

	endpoint, err := h.getEndpointUseCase.Execute(c.Context(), activeOrganizationID, endpointUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get endpoint",
			"errors":  err.Error(),
		})
	}

	return c.JSON(presenter.NewEndpointDetailResponse(endpoint))
}

func (h *EndpointHandler) UpdateEndpoint(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
			"errors":  err.Error(),
		})
	}

	var request endpointDTO.UpdateEndpointInput
	if err := c.Bind().JSON(&request); err != nil {
		fmt.Println("error1", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
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
			"errors":  err.Error(),
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
			"errors":  err.Error(),
		})
	}

	var request endpointDTO.ImportEndpointsInput
	if err := c.Bind().Form(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
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
			"errors":  err.Error(),
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
			"errors":  err.Error(),
		})
	}

	endpointID := c.Params("id")
	endpointUUID, err := uuid.Parse(endpointID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid endpoint ID",
			"errors":  err.Error(),
		})
	}

	hasSteps, err := h.endpointHasStepUseCase.Execute(c.Context(), endpointUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check if endpoint has steps",
			"errors":  err.Error(),
		})
	}

	if hasSteps {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Endpoint has steps",
			"errors":  errors.New("endpoint has steps").Error(),
		})
	}

	err = h.deleteEndpointUseCase.Execute(c.Context(), endpointUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete endpoint",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
