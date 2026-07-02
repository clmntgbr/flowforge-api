package handler

import (
	"flowforge-api/handler/context"
	organizationDTO "flowforge-api/infrastructure/organization"
	"flowforge-api/presenter"
	"flowforge-api/usecase/organization"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	listOrganizationsUseCase    *organization.ListOrganizationsUseCase
	createOrganizationUseCase   *organization.CreateOrganizationUseCase
	getOrganizationByIDUseCase  *organization.GetOrganizationByIDUseCase
	updateOrganizationUseCase   *organization.UpdateOrganizationUseCase
	activateOrganizationUseCase *organization.ActivateOrganizationUseCase
}

func NewOrganizationHandler(
	listOrganizationsUseCase *organization.ListOrganizationsUseCase,
	createOrganizationUseCase *organization.CreateOrganizationUseCase,
	getOrganizationByIDUseCase *organization.GetOrganizationByIDUseCase,
	updateOrganizationUseCase *organization.UpdateOrganizationUseCase,
	activateOrganizationUseCase *organization.ActivateOrganizationUseCase,
) *OrganizationHandler {
	return &OrganizationHandler{
		listOrganizationsUseCase:    listOrganizationsUseCase,
		createOrganizationUseCase:   createOrganizationUseCase,
		getOrganizationByIDUseCase:  getOrganizationByIDUseCase,
		updateOrganizationUseCase:   updateOrganizationUseCase,
		activateOrganizationUseCase: activateOrganizationUseCase,
	}
}

func (h *OrganizationHandler) GetOrganizations(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	organizations, err := h.listOrganizationsUseCase.Execute(c.Context(), user, activeOrganizationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to list organizations",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(presenter.NewOrganizationListResponses(organizations))
}

func (h *OrganizationHandler) CreateOrganization(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	var request organizationDTO.CreateOrganizationInput
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

	_, err = h.createOrganizationUseCase.Execute(c.Context(), user, request.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create organization",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *OrganizationHandler) GetOrganizationByID(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	organizationID := c.Params("id")
	organizationUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid organization ID",
			"errors":  err.Error(),
		})
	}

	organization, err := h.getOrganizationByIDUseCase.Execute(c.Context(), user, organizationUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get organization",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(presenter.NewOrganizationDetailResponse(organization))
}

func (h *OrganizationHandler) UpdateOrganization(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	organizationID := c.Params("id")
	organizationUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid organization ID",
			"errors":  err.Error(),
		})
	}

	var request organizationDTO.UpdateOrganizationInput
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

	organization, err := h.updateOrganizationUseCase.Execute(c.Context(), user, organizationUUID, request.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update organization",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(presenter.NewOrganizationDetailResponse(organization))
}

func (h *OrganizationHandler) ActivateOrganization(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	organizationID := c.Params("id")
	organizationUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid organization ID",
			"errors":  err.Error(),
		})
	}

	_, err = h.activateOrganizationUseCase.Execute(c.Context(), user, organizationUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to activate organization",
			"errors":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
