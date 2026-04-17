package ctxutil

import (
	"forgeflow-api/domain"
	"forgeflow-api/errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	UserKey           = "user"
	OrganizationIDKey = "organization_id"
)

func GetUser(c fiber.Ctx) (*domain.User, error) {
	user, ok := c.Locals(UserKey).(*domain.User)
	if !ok {
		return nil, errors.ErrUserNotAuthenticated
	}
	return user, nil
}

func GetOrganizationID(c fiber.Ctx) (uuid.UUID, error) {
	organizationID, ok := c.Locals(OrganizationIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.ErrOrganizationNotFound
	}
	return organizationID, nil
}

func SetUser(c fiber.Ctx, user domain.User) {
	c.Locals(UserKey, &user)
}

func SetOrganizationID(c fiber.Ctx, organizationID uuid.UUID) {
	c.Locals(OrganizationIDKey, organizationID)
}
