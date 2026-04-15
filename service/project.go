package service

import (
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/repository"
	"forgeflow-api/rules"

	"github.com/gofiber/fiber/v3"
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

func (s *ProjectService) GetProjects(c fiber.Ctx, user *domain.User, activeProjectID uuid.UUID) ([]dto.ProjectOutput, error) {
	projects, err := s.projectRepository.FindAllByUserID(c.Context(), user.ID)
	if err != nil {
		return nil, err
	}

	return dto.NewProjectsOutput(projects, activeProjectID), nil
}

func (s *ProjectService) GetProjectByID(c fiber.Ctx, user *domain.User, projectUUID uuid.UUID) (dto.ProjectOutput, error) {
	project, err := s.projectRepository.FindByUserIDAndProjectID(c.Context(), projectUUID, user.ID)
	if err != nil {
		return dto.ProjectOutput{}, err
	}

	return dto.NewProjectOutput(*project, project.ID), nil
}
