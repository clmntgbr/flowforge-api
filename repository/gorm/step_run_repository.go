package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
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

func (r *stepRunRepository) GetAllByWorkflowRunID(ctx context.Context, workflowRunID uuid.UUID) ([]entity.StepRun, error) {
	var stepRuns []entity.StepRun
	err := dbWithContext(ctx, r.db).
		Where("workflow_run_id = ?", workflowRunID).
		Find(&stepRuns).Error
	if err != nil {
		return nil, err
	}
	return stepRuns, nil
}

func (r *stepRunRepository) CancelRunningByWorkflowRunID(ctx context.Context, workflowRunID uuid.UUID) error {
	var stepRuns []entity.StepRun
	err := dbWithContext(ctx, r.db).
		Where("workflow_run_id = ?", workflowRunID).
		Where("status = ?", enum.StepRunStatusRunning).
		Find(&stepRuns).Error
	if err != nil {
		return err
	}

	for i := range stepRuns {
		stepRuns[i].Status = enum.StepRunStatusCanceled
		stepRuns[i].Statuses = append(stepRuns[i].Statuses, enum.StepRunStatusCanceled)
		if err := dbWithContext(ctx, r.db).Save(&stepRuns[i]).Error; err != nil {
			return err
		}
	}

	return nil
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

func (r *stepRunRepository) GetByStepID(ctx context.Context, stepID uuid.UUID) ([]entity.StepRun, error) {
	var stepRuns []entity.StepRun
	err := dbWithContext(ctx, r.db).
		Where("step_id = ?", stepID).
		Where("status = ?", enum.StepRunStatusCompleted).
		Where("response IS NOT NULL AND response != ''").
		Order("created_at DESC").
		Limit(10).
		Find(&stepRuns).Error
	if err != nil {
		return nil, err
	}
	return stepRuns, nil
}

func (r *stepRunRepository) GetLatestCompletedByStepID(ctx context.Context, stepID uuid.UUID) (*entity.StepRun, error) {
	var stepRun entity.StepRun
	err := dbWithContext(ctx, r.db).
		Where("step_id = ?", stepID).
		Where("status = ?", enum.StepRunStatusCompleted).
		Where("response IS NOT NULL AND response != ''").
		Order("created_at DESC").
		First(&stepRun).Error
	if err != nil {
		return nil, err
	}
	return &stepRun, nil
}
