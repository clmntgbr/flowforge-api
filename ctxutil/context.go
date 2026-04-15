package ctxutil

import (
	"go-api/domain"
	"go-api/errors"

	"github.com/gofiber/fiber/v3"
)

const (
	UserKey = "user"
)

func GetUser(c fiber.Ctx) (*domain.User, error) {
	user, ok := c.Locals(UserKey).(*domain.User)
	if !ok {
		return nil, errors.ErrUserNotAuthenticated
	}
	return user, nil
}

func SetUser(c fiber.Ctx, user domain.User) {
	c.Locals(UserKey, &user)
}
