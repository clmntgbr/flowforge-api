package handler

import (
	"flowforge-api/handler/context"
	organizationDTO "flowforge-api/infrastructure/organization"
	"flowforge-api/presenter"
	"flowforge-api/usecase/organization"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	listOrganizationsUseCase   *organization.ListOrganizationsUseCase
	createOrganizationUseCase  *organization.CreateOrganizationUseCase
	getOrganizationByIDUseCase *organization.GetOrganizationByIDUseCase
}

func NewOrganizationHandler(listOrganizationsUseCase *organization.ListOrganizationsUseCase, createOrganizationUseCase *organization.CreateOrganizationUseCase, getOrganizationByIDUseCase *organization.GetOrganizationByIDUseCase) *OrganizationHandler {
	return &OrganizationHandler{listOrganizationsUseCase: listOrganizationsUseCase, createOrganizationUseCase: createOrganizationUseCase, getOrganizationByIDUseCase: getOrganizationByIDUseCase}
}

func (h *OrganizationHandler) GetOrganizations(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	organizations, err := h.listOrganizationsUseCase.Execute(c.Context(), user, activeOrganizationID)
	if err != nil {
		log.Println("Failed to list organizations: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to list organizations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(presenter.NewOrganizationListResponses(organizations))
}

func (h *OrganizationHandler) CreateOrganization(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var request organizationDTO.CreateOrganizationInput
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

	organization, err := h.createOrganizationUseCase.Execute(c.Context(), user, request.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create organization",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.NewOrganizationDetailResponse(organization))
}

func (h *OrganizationHandler) GetOrganizationByID(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	organizationID := c.Params("id")
	organizationUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid organization ID",
		})
	}

	organization, err := h.getOrganizationByIDUseCase.Execute(c.Context(), user, organizationUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get organization",
		})
	}

	return c.Status(fiber.StatusOK).JSON(presenter.NewOrganizationDetailResponse(organization))
}
