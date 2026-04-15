package ctxutil

import (
	"forgeflow-api/domain"
	"forgeflow-api/errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	UserKey      = "user"
	ProjectIDKey = "project_id"
)

func GetUser(c fiber.Ctx) (*domain.User, error) {
	user, ok := c.Locals(UserKey).(*domain.User)
	if !ok {
		return nil, errors.ErrUserNotAuthenticated
	}
	return user, nil
}

func GetProjectID(c fiber.Ctx) (uuid.UUID, error) {
	projectID, ok := c.Locals(ProjectIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.ErrProjectNotFound
	}
	return projectID, nil
}

func SetUser(c fiber.Ctx, user domain.User) {
	c.Locals(UserKey, &user)
}

func SetProjectID(c fiber.Ctx, projectID uuid.UUID) {
	c.Locals(ProjectIDKey, projectID)
}
