package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stepRepository struct {
	db *gorm.DB
}

func NewStepRepository(db *gorm.DB) repository.StepRepository {
	return &stepRepository{db: db}
}

func (r *stepRepository) Create(ctx context.Context, step *entity.Step) error {
	return r.db.WithContext(ctx).Create(step).Error
}

func (r *stepRepository) Update(ctx context.Context, step *entity.Step) error {
	return r.db.WithContext(ctx).Save(step).Error
}

func (r *stepRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Step{}, id).Error
}

func (r *stepRepository) GetByIDAndOrganizationIDAndWorkflowID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, id uuid.UUID) (entity.Step, error) {
	var step entity.Step

	err := r.db.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = steps.workflow_id").
		Where("steps.id = ? AND workflows.organization_id = ? AND workflows.id = ?", id, organizationID, workflowID).
		First(&step).Error
	if err != nil {
		return entity.Step{}, err
	}
	return step, nil
}
