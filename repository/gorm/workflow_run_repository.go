package gorm

import (
	"context"
	"flowforge-api/domain/entity"
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
	return r.db.WithContext(ctx).Create(workflowRun).Error
}

func (r *workflowRunRepository) Update(ctx context.Context, workflowRun *entity.WorkflowRun) error {
	return r.db.WithContext(ctx).Save(workflowRun).Error
}

func (r *workflowRunRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Workflow{}, id).Error
}

func (r *workflowRunRepository) GetByWorkflowID(ctx context.Context, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.WorkflowRun, int64, error) {
	var workflowRuns []entity.WorkflowRun

	db := r.db.WithContext(ctx).Model(&entity.WorkflowRun{}).
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
