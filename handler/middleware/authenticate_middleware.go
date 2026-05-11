package middleware

import (
	"flowforge-api/presenter"
	"flowforge-api/usecase/auth"
	"flowforge-api/usecase/clerk"
	"flowforge-api/usecase/organization"
	"flowforge-api/usecase/user"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthenticateMiddleware struct {
	validateTokenUseCase      *auth.ValidateTokenUseCase
	getUserUseCase            *clerk.GetUserUseCase
	createUserUseCase         *user.CreateUserUseCase
	createOrganizationUseCase *organization.CreateOrganizationUseCase
	updateUserUseCase         *user.UpdateUserUseCase
}

func NewAuthenticateMiddleware(validateTokenUseCase *auth.ValidateTokenUseCase, getUserUseCase *clerk.GetUserUseCase, createUserUseCase *user.CreateUserUseCase, createOrganizationUseCase *organization.CreateOrganizationUseCase, updateUserUseCase *user.UpdateUserUseCase) *AuthenticateMiddleware {
	return &AuthenticateMiddleware{
		validateTokenUseCase:      validateTokenUseCase,
		getUserUseCase:            getUserUseCase,
		createUserUseCase:         createUserUseCase,
		createOrganizationUseCase: createOrganizationUseCase,
		updateUserUseCase:         updateUserUseCase,
	}
}

func (m *AuthenticateMiddleware) Protected() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid authorization header format",
			})
		}

		if parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization scheme must be Bearer",
			})
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token cannot be empty",
			})
		}

		output, err := m.validateTokenUseCase.Execute(c.Context(), presenter.ValidateTokenInput{
			Token: tokenString,
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		if output.User == nil {
			clerkUser, err := m.getUserUseCase.Execute(c.Context(), output.Claims.Subject)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Failed to get user",
				})
			}

			user, err := m.createUserUseCase.Execute(c.Context(), output.Claims.Subject, clerkUser.FirstName, clerkUser.LastName, clerkUser.Banned)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Failed to create user",
				})
			}

			organization, err := m.createOrganizationUseCase.Execute(c.Context(), user, "Default Organization")
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Failed to create organization",
				})
			}

			organizationID, err := uuid.Parse(organization.ID)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Failed to create organization",
				})
			}

			user.ActiveOrganizationID = &organizationID
			if err := m.updateUserUseCase.Execute(c.Context(), user); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Failed to update user",
				})
			}

			output.User = user
		}

		if output.User.Banned {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User is banned",
			})
		}

		return c.Next()
	}
}
