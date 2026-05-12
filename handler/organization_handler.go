package handler

import (
	"flowforge-api/handler/context"
	"flowforge-api/presenter"
	"flowforge-api/usecase/organization"
	"log"

	"github.com/gofiber/fiber/v3"
)

type OrganizationHandler struct {
	listOrganizationsUseCase *organization.ListOrganizationsUseCase
}

func NewOrganizationHandler(listOrganizationsUseCase *organization.ListOrganizationsUseCase) *OrganizationHandler {
	return &OrganizationHandler{listOrganizationsUseCase: listOrganizationsUseCase}
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
