package main

import (
	"flowforge-api/handler"
	"flowforge-api/handler/middleware"
	infraClerk "flowforge-api/infrastructure/clerk"
	"flowforge-api/infrastructure/config"
	repoGorm "flowforge-api/repository/gorm"
	"flowforge-api/usecase/auth"
	"flowforge-api/usecase/clerk"
	"flowforge-api/usecase/organization"
	"flowforge-api/usecase/user"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkMiddleware        *middleware.ClerkMiddleware
	ClerkHandler           *handler.ClerkHandler
	UserHandler            *handler.UserHandler
	OrganizationHandler    *handler.OrganizationHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	jwksProvider, err := infraClerk.NewJWKSProvider(env)
	if err != nil {
		log.Fatalf("failed to create JWKS provider: %v", err)
	}

	userRepo := repoGorm.NewUserRepository(db)
	organizationRepo := repoGorm.NewOrganizationRepository(db)

	validateTokenUseCase := auth.NewValidateTokenUseCase(jwksProvider, userRepo)
	fetchUserUseCase := clerk.NewFetchUserUseCase(env)
	getUserByClerkIDUseCase := user.NewGetUserByClerkIDUseCase(userRepo)
	createUserUseCase := user.NewCreateUserUseCase(userRepo)
	createOrganizationUseCase := organization.NewCreateOrganizationUseCase(organizationRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo)
	deleteUserByClerkIDUseCase := user.NewDeleteUserByClerkIDUseCase(userRepo)
	listOrganizationsUseCase := organization.NewListOrganizationsUseCase(organizationRepo)
	getOrganizationByIDUseCase := organization.NewGetOrganizationByIDUseCase(organizationRepo)
	updateOrganizationUseCase := organization.NewUpdateOrganizationUseCase(organizationRepo)
	activateOrganizationUseCase := organization.NewActivateOrganizationUseCase(organizationRepo)

	clerkMiddleware := middleware.NewClerkMiddleware(
		env.ClerkWebhookSecret,
	)

	authenticateMiddleware := middleware.NewAuthenticateMiddleware(
		validateTokenUseCase,
		fetchUserUseCase,
		createUserUseCase,
		createOrganizationUseCase,
		updateUserUseCase,
	)

	clerkHandler := handler.NewClerkHandler(
		getUserByClerkIDUseCase,
		createUserUseCase,
		createOrganizationUseCase,
		updateUserUseCase,
		deleteUserByClerkIDUseCase,
	)

	userHandler := handler.NewUserHandler()

	organizationHandler := handler.NewOrganizationHandler(
		listOrganizationsUseCase,
		createOrganizationUseCase,
		getOrganizationByIDUseCase,
		updateOrganizationUseCase,
		activateOrganizationUseCase,
	)

	return &Container{
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkMiddleware:        clerkMiddleware,
		ClerkHandler:           clerkHandler,
		UserHandler:            userHandler,
		OrganizationHandler:    organizationHandler,
	}
}
