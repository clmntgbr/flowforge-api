package service

import (
	"forgeflow-api/domain"
	"forgeflow-api/repository"
	"forgeflow-api/rules"

	"github.com/gofiber/fiber/v3"
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

func (s *ProjectService) CreateProject(c fiber.Ctx, user *domain.User, name string) (*domain.Project, error) {

	if err := s.projectRules.MaxProjectsPerUser(c.Context(), user.ID); err != nil {
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
