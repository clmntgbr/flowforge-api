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
	OrganizationService *service.OrganizationService

	WebhookClerkHandler *handler.WebhookClerkHandler
	UserHandler         *handler.UserHandler
	OrganizationHandler *handler.OrganizationHandler

	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkWebhookMiddleware *middleware.ClerkWebhookMiddleware
}

func New(db *gorm.DB, cfg *config.Config) *Dependencies {
	userRepo := repository.NewUserRepository(db)
	organizationRepo := repository.NewOrganizationRepository(db)

	organizationRules := rules.NewOrganizationRules(organizationRepo)

	organizationService := service.NewOrganizationService(organizationRepo, organizationRules)
	authenticateService := service.NewAuthenticateService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	clerkService := service.NewClerkService(cfg)

	webhookClerkHandler := handler.NewWebhookClerkHandler(userService, organizationService, userRepo)
	userHandler := handler.NewUserHandler(userService)
	organizationHandler := handler.NewOrganizationHandler(organizationService)

	clerkWebhookMiddleware := middleware.NewClerkWebhookMiddleware(cfg.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(authenticateService, clerkService, userService, organizationService, userRepo)

	return &Dependencies{
		UserRepo:               userRepo,
		AuthenticateService:    authenticateService,
		UserService:            userService,
		OrganizationService:    organizationService,
		WebhookClerkHandler:    webhookClerkHandler,
		UserHandler:            userHandler,
		OrganizationHandler:    organizationHandler,
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkWebhookMiddleware: clerkWebhookMiddleware,
		ClerkService:           clerkService,
	}
}
