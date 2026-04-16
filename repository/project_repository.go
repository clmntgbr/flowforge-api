package repository

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *domain.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) Update(project *domain.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(project *domain.Project) error {
	return r.db.Delete(project).Error
}

func (r *ProjectRepository) CountProjectsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Project{}).
		Joins("JOIN user_projects ON user_projects.project_id = projects.id").
		Where("user_projects.user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProjectRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Project, error) {
	var projects []domain.Project

	db := r.db.WithContext(ctx).
		Model(&domain.Project{}).
		Joins("JOIN user_projects ON user_projects.project_id = projects.id").
		Where("user_projects.user_id = ?", userID)

	err := db.Find(&projects).Error
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) FindByUserIDAndProjectID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*domain.Project, error) {
	var project domain.Project

	err := r.db.WithContext(ctx).
		Joins("JOIN user_projects ON user_projects.project_id = projects.id").
		Where("projects.id = ? AND user_projects.user_id = ?", projectID, userID).
		First(&project).Error

	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) ActivateProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) (*domain.Project, error) {
	var project domain.Project
	err := r.db.WithContext(ctx).
		Joins("JOIN user_projects ON user_projects.project_id = projects.id").
		Where("projects.id = ? AND user_projects.user_id = ?", projectID, userID).
		First(&project).Error

	if err != nil {
		return nil, errors.ErrProjectNotFound
	}

	err = r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("active_project_id", projectID).Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}
