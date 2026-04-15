package ctxutil

import (
	"forgeflow-api/domain"
	"forgeflow-api/errors"

	"github.com/gofiber/fiber/v3"
)

const (
	UserKey    = "user"
	ProjectKey = "project"
)

func GetUser(c fiber.Ctx) (*domain.User, error) {
	user, ok := c.Locals(UserKey).(*domain.User)
	if !ok {
		return nil, errors.ErrUserNotAuthenticated
	}
	return user, nil
}

func GetProject(c fiber.Ctx) (*domain.Project, error) {
	project, ok := c.Locals(ProjectKey).(*domain.Project)
	if !ok {
		return nil, errors.ErrProjectNotFound
	}
	return project, nil
}

func SetUser(c fiber.Ctx, user domain.User) {
	c.Locals(UserKey, &user)
}

func SetProject(c fiber.Ctx, project domain.Project) {
	c.Locals(ProjectKey, &project)
}
