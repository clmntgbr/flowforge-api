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
	EndpointService     *service.EndpointService
	WorkflowService     *service.WorkflowService

	WebhookClerkHandler *handler.WebhookClerkHandler
	UserHandler         *handler.UserHandler
	OrganizationHandler *handler.OrganizationHandler
	EndpointHandler     *handler.EndpointHandler
	WorkflowHandler     *handler.WorkflowHandler

	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkWebhookMiddleware *middleware.ClerkWebhookMiddleware
}

func New(db *gorm.DB, cfg *config.Config) *Dependencies {
	userRepo := repository.NewUserRepository(db)
	organizationRepo := repository.NewOrganizationRepository(db)
	endpointRepo := repository.NewEndpointRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)

	organizationRules := rules.NewOrganizationRules(organizationRepo)

	organizationService := service.NewOrganizationService(organizationRepo, organizationRules)
	authenticateService := service.NewAuthenticateService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	clerkService := service.NewClerkService(cfg)
	endpointService := service.NewEndpointService(endpointRepo)
	workflowService := service.NewWorkflowService(workflowRepo)

	webhookClerkHandler := handler.NewWebhookClerkHandler(userService, organizationService, userRepo)
	userHandler := handler.NewUserHandler(userService)
	organizationHandler := handler.NewOrganizationHandler(organizationService)
	endpointHandler := handler.NewEndpointHandler(endpointService)
	workflowHandler := handler.NewWorkflowHandler(workflowService)

	clerkWebhookMiddleware := middleware.NewClerkWebhookMiddleware(cfg.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(authenticateService, clerkService, userService, organizationService, userRepo)

	return &Dependencies{
		UserRepo:               userRepo,
		AuthenticateService:    authenticateService,
		UserService:            userService,
		OrganizationService:    organizationService,
		EndpointService:        endpointService,
		WorkflowService:        workflowService,
		WebhookClerkHandler:    webhookClerkHandler,
		UserHandler:            userHandler,
		OrganizationHandler:    organizationHandler,
		EndpointHandler:        endpointHandler,
		WorkflowHandler:        workflowHandler,
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkWebhookMiddleware: clerkWebhookMiddleware,
		ClerkService:           clerkService,
	}
}
