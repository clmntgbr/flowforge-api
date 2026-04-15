// deps/deps.go
package deps

import (
	"forgeflow-api/config"
	"forgeflow-api/handler"
	"forgeflow-api/middleware"
	"forgeflow-api/repository"
	"forgeflow-api/rules"
	"forgeflow-api/service"

	"gorm.io/gorm"
)

type Dependencies struct {
	UserRepo *repository.UserRepository

	AuthenticateService *service.AuthenticateService
	ClerkService        *service.ClerkService
	UserService         *service.UserService
	ProjectService      *service.ProjectService

	WebhookClerkHandler *handler.WebhookClerkHandler
	UserHandler         *handler.UserHandler
	ProjectHandler      *handler.ProjectHandler

	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkWebhookMiddleware *middleware.ClerkWebhookMiddleware
}

func New(db *gorm.DB, cfg *config.Config) *Dependencies {
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	projectRules := rules.NewProjectRules(projectRepo)

	projectService := service.NewProjectService(projectRepo, projectRules)
	authenticateService := service.NewAuthenticateService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	clerkService := service.NewClerkService(cfg)

	webhookClerkHandler := handler.NewWebhookClerkHandler(userService, projectService, userRepo)
	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService)

	clerkWebhookMiddleware := middleware.NewClerkWebhookMiddleware(cfg.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(authenticateService, clerkService, userService, projectService, userRepo)

	return &Dependencies{
		UserRepo:               userRepo,
		AuthenticateService:    authenticateService,
		UserService:            userService,
		ProjectService:         projectService,
		WebhookClerkHandler:    webhookClerkHandler,
		UserHandler:            userHandler,
		ProjectHandler:         projectHandler,
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkWebhookMiddleware: clerkWebhookMiddleware,
		ClerkService:           clerkService,
	}
}
