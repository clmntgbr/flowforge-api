// deps/deps.go
package deps

import (
	"go-api/config"
	"go-api/handler"
	"go-api/middleware"
	"go-api/repository"
	"go-api/service"

	"gorm.io/gorm"
)

type Dependencies struct {
	UserRepo *repository.UserRepository

	AuthenticateService *service.AuthenticateService
	WebhookClerkService *service.WebhookClerkService
	ClerkService        *service.ClerkService
	UserService         *service.UserService

	WebhookClerkHandler *handler.WebhookClerkHandler
	UserHandler         *handler.UserHandler

	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkWebhookMiddleware *middleware.ClerkWebhookMiddleware
}

func New(db *gorm.DB, cfg *config.Config) *Dependencies {
	userRepo := repository.NewUserRepository(db)

	authenticateService := service.NewAuthenticateService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	webhookClerkService := service.NewWebhookClerkService(userRepo, userService)
	clerkService := service.NewClerkService(cfg)

	webhookClerkHandler := handler.NewWebhookClerkHandler(webhookClerkService)
	userHandler := handler.NewUserHandler(userService)

	clerkWebhookMiddleware := middleware.NewClerkWebhookMiddleware(cfg.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(authenticateService, clerkService, userService, userRepo)

	return &Dependencies{
		UserRepo:               userRepo,
		AuthenticateService:    authenticateService,
		WebhookClerkService:    webhookClerkService,
		UserService:            userService,
		WebhookClerkHandler:    webhookClerkHandler,
		UserHandler:            userHandler,
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkWebhookMiddleware: clerkWebhookMiddleware,
		ClerkService:           clerkService,
	}
}
