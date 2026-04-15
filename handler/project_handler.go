package handler

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/service"

	"github.com/gofiber/fiber/v3"
)

type ProjectHandler struct {
	BaseHandler
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) GetProjects(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	activeProjectID, err := ctxutil.GetProjectID(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	output, err := h.projectService.GetProjects(c, user, activeProjectID)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
