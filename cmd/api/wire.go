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

type Dependencies struct {
	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkMiddleware        *middleware.ClerkMiddleware
	UserHandler            *handler.UserHandler
}

func NewWire(db *gorm.DB, env *config.Config) *Dependencies {
	userRepo := repoGorm.NewUserRepository(db)
	organizationRepo := repoGorm.NewOrganizationRepository(db)

	jwksProvider, err := infraClerk.NewJWKSProvider(env)
	if err != nil {
		log.Fatalf("failed to create JWKS provider: %v", err)
	}

	validateTokenUseCase := auth.NewValidateTokenUseCase(jwksProvider, userRepo)
	fetchUserUseCase := clerk.NewFetchUserUseCase(env)
	createUserUseCase := user.NewCreateUserUseCase(userRepo)
	createOrganizationUseCase := organization.NewCreateOrganizationUseCase(organizationRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo)

	clerkMiddleware := middleware.NewClerkMiddleware(env.ClerkSecretKey)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(validateTokenUseCase, fetchUserUseCase, createUserUseCase, createOrganizationUseCase, updateUserUseCase)

	userHandler := handler.NewUserHandler()

	return &Dependencies{
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkMiddleware:        clerkMiddleware,
		UserHandler:            userHandler,
	}
}
