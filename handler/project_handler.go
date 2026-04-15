package handler

import (
	"fmt"
	"forgeflow-api/ctxutil"
	"forgeflow-api/usecase"

	"github.com/gofiber/fiber/v3"
)

type ProjectHandler struct {
	BaseHandler
	getProjectsUsecase *usecase.GetProjectsUsecase
}

func NewProjectHandler(getProjectsUsecase *usecase.GetProjectsUsecase) *ProjectHandler {
	return &ProjectHandler{
		getProjectsUsecase: getProjectsUsecase,
	}
}

func (h *ProjectHandler) GetProjects(c fiber.Ctx) error {
	user, err := ctxutil.GetUser(c)
	if err != nil {
		return h.sendNotFound(c, err)
	}

	fmt.Println("user", user)

	activeProject, err := ctxutil.GetProject(c)
	if err != nil {
		return h.sendNotFound(c, err)
	}

	fmt.Println("activeProject", activeProject)

	output, err := h.getProjectsUsecase.GetProjectsByUserID(c.Context(), user, activeProject)
	if err != nil {
		return h.sendInternalError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

func (h *ProjectHandler) GetProjectByID(c fiber.Ctx) error {
	return nil
}

func (h *ProjectHandler) CreateProject(c fiber.Ctx) error {
	return nil
}

func (h *ProjectHandler) UpdateProject(c fiber.Ctx) error {
	return nil
}

func (h *ProjectHandler) ActivateProject(c fiber.Ctx) error {
	return nil
}
