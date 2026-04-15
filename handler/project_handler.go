package handler

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/dto"
	"forgeflow-api/errors"
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

func (h *ProjectHandler) GetProjectByID(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	projectUUID, err := h.parseUUIDParam(c, "id", errors.ErrInvalidProjectID)
	if err != nil {
		return err
	}

	project, err := h.projectService.GetProjectByID(c, user, projectUUID)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(project)
}

func (h *ProjectHandler) CreateProject(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendUnauthorized(c)
	}

	var req dto.CreateProjectInput
	if err := h.bindAndValidate(c, &req); err != nil {
		return nil
	}

	project, err := h.projectService.CreateProject(c, user, req.Name)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}
