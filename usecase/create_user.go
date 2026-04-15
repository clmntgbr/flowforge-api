package usecase

import (
	"forgeflow-api/domain"
	"forgeflow-api/service"

	"github.com/gofiber/fiber/v3"
)

type CreateUserUsecase struct {
	userService    *service.UserService
	projectService *service.ProjectService
}

func NewCreateUserUsecase(userService *service.UserService, projectService *service.ProjectService) *CreateUserUsecase {
	return &CreateUserUsecase{
		userService:    userService,
		projectService: projectService,
	}
}

func (u *CreateUserUsecase) CreateUser(c fiber.Ctx, clerkID string, firstName string, lastName string, banned bool) (*domain.User, *domain.Project, error) {
	user, err := u.userService.CreateUser(c, clerkID, firstName, lastName, banned)
	if err != nil {
		return nil, nil, err
	}

	project, err := u.projectService.CreateProject(c, user, "Default Project")
	if err != nil {
		return nil, nil, err
	}

	return user, project, nil
}
