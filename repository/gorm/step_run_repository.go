package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stepRunRepository struct {
	db *gorm.DB
}

func NewStepRunRepository(db *gorm.DB) repository.StepRunRepository {
	return &stepRunRepository{db: db}
}

func (r *stepRunRepository) Create(ctx context.Context, stepRun *entity.StepRun) error {
	return dbWithContext(ctx, r.db).Create(stepRun).Error
}

func (r *stepRunRepository) Update(ctx context.Context, stepRun *entity.StepRun) error {
	return dbWithContext(ctx, r.db).Save(stepRun).Error
}

func (r *stepRunRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.StepRun{}, id).Error
}

func (r *stepRunRepository) GetByWorkflowRunID(ctx context.Context, workflowRunID uuid.UUID) (*entity.StepRun, error) {
	var stepRun entity.StepRun
	err := dbWithContext(ctx, r.db).
		Where("workflow_run_id = ?", workflowRunID).
		First(&stepRun).Error
	if err != nil {
		return nil, err
	}
	return &stepRun, nil
}

func (r *stepRunRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.StepRun, error) {
	var stepRun entity.StepRun
	err := dbWithContext(ctx, r.db).
		Where("id = ?", id).
		First(&stepRun).Error
	if err != nil {
		return nil, err
	}
	return &stepRun, nil
}
