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
	WebhookClerkHandler *handler.WebhookClerkHandler
	UserHandler         *handler.UserHandler
	OrganizationHandler *handler.OrganizationHandler
	EndpointHandler     *handler.EndpointHandler
	WorkflowHandler     *handler.WorkflowHandler
	ConnexionHandler    *handler.ConnexionHandler
	StepHandler         *handler.StepHandler

	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkWebhookMiddleware *middleware.ClerkWebhookMiddleware
}

func New(db *gorm.DB, cfg *config.Config) *Dependencies {
	userRepo := repository.NewUserRepository(db)
	organizationRepo := repository.NewOrganizationRepository(db)
	endpointRepo := repository.NewEndpointRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	stepRepo := repository.NewStepRepository(db)
	connexionRepo := repository.NewConnexionRepository(db)

	organizationRules := rules.NewOrganizationRules(organizationRepo)

	organizationService := service.NewOrganizationService(organizationRepo, organizationRules)
	authenticateService := service.NewAuthenticateService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	clerkService := service.NewClerkService(cfg)
	endpointService := service.NewEndpointService(endpointRepo)
	workflowService := service.NewWorkflowService(workflowRepo)
	stepService := service.NewStepService(stepRepo, endpointRepo)
	connexionService := service.NewConnexionService(connexionRepo)

	webhookClerkHandler := handler.NewWebhookClerkHandler(userService, organizationService, userRepo)
	userHandler := handler.NewUserHandler(userService)
	organizationHandler := handler.NewOrganizationHandler(organizationService)
	endpointHandler := handler.NewEndpointHandler(endpointService)
	workflowHandler := handler.NewWorkflowHandler(workflowService, stepService)
	connexionHandler := handler.NewConnexionHandler(connexionService, workflowService)
	stepHandler := handler.NewStepHandler(stepService)

	clerkWebhookMiddleware := middleware.NewClerkWebhookMiddleware(cfg.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(authenticateService, clerkService, userService, organizationService, userRepo)

	return &Dependencies{
		WebhookClerkHandler:    webhookClerkHandler,
		UserHandler:            userHandler,
		OrganizationHandler:    organizationHandler,
		EndpointHandler:        endpointHandler,
		WorkflowHandler:        workflowHandler,
		ConnexionHandler:       connexionHandler,
		StepHandler:            stepHandler,
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkWebhookMiddleware: clerkWebhookMiddleware,
	}
}
