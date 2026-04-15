package usecase

import (
	"context"
	"forgeflow-api/domain"
)

type CreateProjectUsecase struct {
	projects ProjectProvisioner
}

func NewCreateProjectUsecase(projects ProjectProvisioner) *CreateProjectUsecase {
	return &CreateProjectUsecase{projects: projects}
}

func (u *CreateProjectUsecase) CreateProject(ctx context.Context, user *domain.User, name string) (*domain.Project, error) {
	return u.projects.CreateProject(ctx, user, name)
}
