package rules

import (
	"context"
	"forgeflow-api/errors"
	"forgeflow-api/repository"

	"github.com/google/uuid"
)

type ProjectRules struct {
	projectRepo *repository.ProjectRepository
}

func NewProjectRules(projectRepo *repository.ProjectRepository) *ProjectRules {
	return &ProjectRules{
		projectRepo: projectRepo,
	}
}
func (p *ProjectRules) MaxProjectsPerUser(ctx context.Context, userID uuid.UUID) error {
	count, err := p.projectRepo.CountProjectsByUserID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	if count >= 10 {
		return errors.ErrMaxProjectsReached
	}

	return nil
}
