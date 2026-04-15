// deps/deps.go
package deps

import (
	"forgeflow-api/config"
	"forgeflow-api/handler"
	"forgeflow-api/middleware"
	"forgeflow-api/repository"
	"forgeflow-api/rules"
	"forgeflow-api/service"
	"forgeflow-api/usecase"

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

	createProjectUsecase := usecase.NewCreateProjectUsecase(projectService)
	createUserUsecase := usecase.NewCreateUserUsecase(userService, createProjectUsecase, userRepo)
	updateUserUsecase := usecase.NewUpdateUserUsecase(userService)
	deleteUserUsecase := usecase.NewDeleteUserUsecase(userService)
	getUserUsecase := usecase.NewGetUserUsecase(userService)
	getProjectsUsecase := usecase.NewGetProjectsUsecase(projectService)

	webhookClerkHandler := handler.NewWebhookClerkHandler(createUserUsecase, updateUserUsecase, deleteUserUsecase)
	userHandler := handler.NewUserHandler(getUserUsecase)
	projectHandler := handler.NewProjectHandler(getProjectsUsecase)

	clerkWebhookMiddleware := middleware.NewClerkWebhookMiddleware(cfg.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(authenticateService, clerkService, createUserUsecase, userRepo)

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
