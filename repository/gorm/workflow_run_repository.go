package gorm

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type workflowRunRepository struct {
	db *gorm.DB
}

func NewWorkflowRunRepository(db *gorm.DB) repository.WorkflowRunRepository {
	return &workflowRunRepository{db: db}
}

func (r *workflowRunRepository) Create(ctx context.Context, workflowRun *entity.WorkflowRun) error {
	return dbWithContext(ctx, r.db).Create(workflowRun).Error
}

func (r *workflowRunRepository) Update(ctx context.Context, workflowRun *entity.WorkflowRun) error {
	return dbWithContext(ctx, r.db).Save(workflowRun).Error
}

func (r *workflowRunRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Workflow{}, id).Error
}

func (r *workflowRunRepository) GetByWorkflowID(ctx context.Context, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.WorkflowRun, int64, error) {
	var workflowRuns []entity.WorkflowRun

	db := dbWithContext(ctx, r.db).Model(&entity.WorkflowRun{}).
		Where("workflow_id = ?", workflowID)

	db, total, err := Paginate(db, query)
	if err != nil {
		return nil, 0, err
	}

	err = db.Preload("StepsRuns").
		Preload("StepsRuns.Step").
		Preload("StepsRuns.Step.Endpoint").
		Preload("StepsRuns.Insight").
		Find(&workflowRuns).Error

	if err != nil {
		return nil, 0, err
	}

	return workflowRuns, total, nil
}

func (r *workflowRunRepository) GetByWorkflowIDAndNotEnded(ctx context.Context, workflowID uuid.UUID) (*entity.WorkflowRun, error) {
	var workflowRun entity.WorkflowRun

	err := dbWithContext(ctx, r.db).
		Where("workflow_id = ?",
			workflowID,
		).
		Where("status != ?", enum.WorkflowRunStatusCompleted).
		Where("status != ?", enum.WorkflowRunStatusFailed).
		Preload("Workflow").
		First(&workflowRun).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &workflowRun, nil
}
