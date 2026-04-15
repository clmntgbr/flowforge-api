package service

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/repository"
	"forgeflow-api/rules"

	"github.com/google/uuid"
)

type ProjectService struct {
	projectRepository *repository.ProjectRepository
	projectRules      *rules.ProjectRules
}

func NewProjectService(projectRepository *repository.ProjectRepository, projectRules *rules.ProjectRules) *ProjectService {
	return &ProjectService{
		projectRepository: projectRepository,
		projectRules:      projectRules,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, user *domain.User, name string) (*domain.Project, error) {

	if err := s.projectRules.MaxProjectsPerUser(ctx, user.ID); err != nil {
		return nil, err
	}

	project := &domain.Project{
		Name: name,
		Users: []domain.User{
			{
				ID: user.ID,
			},
		},
	}

	if err := s.projectRepository.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Project, error) {
	return s.projectRepository.FindAllByUserID(ctx, userID)
}
