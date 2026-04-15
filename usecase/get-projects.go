package usecase

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
)

type GetProjectsUsecase struct {
	projects ProjectByUserReader
}

func NewGetProjectsUsecase(projects ProjectByUserReader) *GetProjectsUsecase {
	return &GetProjectsUsecase{projects: projects}
}

func (u *GetProjectsUsecase) GetProjectsByUserID(ctx context.Context, user *domain.User, activeProject *domain.Project) ([]dto.ProjectOutput, error) {
	projects, err := u.projects.GetProjectsByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return dto.NewProjectsOutput(projects, activeProject.ID), nil
}
