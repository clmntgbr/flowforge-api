package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type workflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) repository.WorkflowRepository {
	return &workflowRepository{db: db}
}

func (r *workflowRepository) Create(ctx context.Context, workflow *entity.Workflow) error {
	return r.db.WithContext(ctx).Create(workflow).Error
}

func (r *workflowRepository) Update(ctx context.Context, workflow *entity.Workflow) error {
	return r.db.WithContext(ctx).Save(workflow).Error
}

func (r *workflowRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Workflow{}, id).Error
}

func (r *workflowRepository) List(ctx context.Context, organizationID uuid.UUID, query paginate.PaginateQuery) ([]entity.Workflow, int64, error) {
	var workflows []entity.Workflow

	db := r.db.WithContext(ctx).Model(&entity.Workflow{}).
		Where("organization_id = ?", organizationID)

	if query.Search != "" {
		db = db.Where("name ILIKE ?", "%"+query.Search+"%")
	}

	db, total, err := Paginate(db, query)
	if err != nil {
		return nil, 0, err
	}

	err = db.Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

func (r *workflowRepository) GetByIDAndOrganizationID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (entity.Workflow, error) {
	var workflow entity.Workflow

	err := r.db.WithContext(ctx).Model(&entity.Workflow{}).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("index ASC")
		}).
		Preload("Steps.Endpoint").
		Preload("Connexions").
		Where("organization_id = ? AND id = ?", organizationID, workflowID).
		First(&workflow).Error

	if err != nil {
		return entity.Workflow{}, err
	}

	return workflow, nil
}

func (r *workflowRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *workflowRepository) GetWorkflowsForExecution(ctx context.Context) ([]entity.Workflow, error) {
	var workflows []entity.Workflow

	err := r.db.WithContext(ctx).
		Where("status = ?", enum.WorkflowStatusActive).
		Where("EXISTS (SELECT 1 FROM steps WHERE steps.workflow_id = workflows.id)").
		Find(&workflows).Error
	if err != nil {
		return nil, err
	}
	return workflows, nil
}
