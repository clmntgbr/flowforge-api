package gorm

import (
	"context"
	"flowforge-api/domain/entity"
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
		Where("organization_id = ? AND id = ?", organizationID, workflowID).
		First(&workflow).Error

	if err != nil {
		return entity.Workflow{}, err
	}

	return workflow, nil
}
