package handler

import (
	"flowforge-api/handler/context"
	"flowforge-api/presenter"
	"flowforge-api/usecase/tag"

	"github.com/gofiber/fiber/v3"
)

type TagHandler struct {
	getTagsUseCase *tag.GetTagsUseCase
}

func NewTagHandler(
	getTagsUseCase *tag.GetTagsUseCase,
) *TagHandler {
	return &TagHandler{
		getTagsUseCase: getTagsUseCase,
	}
}

func (h *TagHandler) GetTags(c fiber.Ctx) error {
	activeOrganizationID, err := context.GetOrganizationID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	tags, err := h.getTagsUseCase.Execute(c.Context(), activeOrganizationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	return c.JSON(presenter.NewTagResponse(tags))
}
