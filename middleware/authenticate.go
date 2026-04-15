package middleware

import (
	"forgeflow-api/ctxutil"
	"forgeflow-api/errors"
	"forgeflow-api/service"
	"strings"

	"forgeflow-api/repository"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthenticateMiddleware struct {
	authenticateService *service.AuthenticateService
	clerkService        *service.ClerkService
	userService         *service.UserService
	projectService      *service.ProjectService
	userRepo            *repository.UserRepository
}

func NewAuthenticateMiddleware(authService *service.AuthenticateService, clerkService *service.ClerkService, userService *service.UserService, projectService *service.ProjectService, userRepo *repository.UserRepository) *AuthenticateMiddleware {
	return &AuthenticateMiddleware{
		authenticateService: authService,
		clerkService:        clerkService,
		userService:         userService,
		projectService:      projectService,
		userRepo:            userRepo,
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

		claims, err := m.authenticateService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": errors.ErrInvalidToken,
			})
		}

		user, err := m.userRepo.FindByClerkID(claims.Subject)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": errors.ErrUserNotFound,
			})
		}

		if user == nil {
			clerkUser, err := m.clerkService.GetUser(claims.Subject)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": errors.ErrClerkUserNotFound,
				})
			}

			user, err = m.userService.CreateUser(c, claims.Subject, clerkUser.FirstName, clerkUser.LastName, clerkUser.Banned)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": errors.ErrUserFailedToCreate,
				})
			}

			project, err := m.projectService.CreateProject(c, user, "Default Project")
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": errors.ErrProjectFailedToCreate,
				})
			}

			projectID, err := uuid.Parse(project.ID)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": errors.ErrProjectFailedToCreate,
				})
			}

			user.ActiveProjectID = &projectID
			if err := m.userRepo.Update(user); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": errors.ErrUserFailedToCreate,
				})
			}
		}

		if user.Banned {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": errors.ErrUserBanned,
			})
		}

		ctxutil.SetUser(c, *user)
		ctxutil.SetProjectID(c, *user.ActiveProjectID)

		return c.Next()
	}
}
