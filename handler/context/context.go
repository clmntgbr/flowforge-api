package context

import (
	"errors"
	"flowforge-api/domain/entity"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	UserKey           = "user"
	OrganizationIDKey = "organization_id"
)

func GetUser(c fiber.Ctx) (*entity.User, error) {
	user, ok := c.Locals(UserKey).(*entity.User)
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func GetOrganizationID(c fiber.Ctx) (uuid.UUID, error) {
	organizationID, ok := c.Locals(OrganizationIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("organization not found")
	}
	return organizationID, nil
}

func SetUser(c fiber.Ctx, user entity.User) {
	c.Locals(UserKey, &user)
}

func SetOrganizationID(c fiber.Ctx, organizationID uuid.UUID) {
	c.Locals(OrganizationIDKey, organizationID)
}
