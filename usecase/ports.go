package usecase

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"

	"github.com/google/uuid"
)

type UserProvisioner interface {
	FindByClerkID(clerkID string) (*domain.User, error)
	CreateUser(id string, firstName string, lastName string, banned bool) (*domain.User, error)
}

type ProjectProvisioner interface {
	CreateProject(ctx context.Context, user *domain.User, name string) (*domain.Project, error)
}

type UserUpdater interface {
	UpdateUser(id string, firstName string, lastName string, banned bool) error
}

type UserDeleter interface {
	DeleteUser(id string) error
}

type UserPresenter interface {
	GetUser(user *domain.User) (*dto.UserOutput, error)
}

type ProjectByUserReader interface {
	GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Project, error)
}
