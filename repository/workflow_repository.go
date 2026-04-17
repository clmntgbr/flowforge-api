package repository

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) Create(workflow *domain.Workflow) error {
	return r.db.Create(workflow).Error
}

func (r *WorkflowRepository) Update(workflow *domain.Workflow) error {
	return r.db.Save(workflow).Error
}

func (r *WorkflowRepository) Delete(workflow *domain.Workflow) error {
	return r.db.Delete(workflow).Error
}

func (r *WorkflowRepository) FindAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, q dto.PaginateQuery) ([]domain.Workflow, int64, error) {
	var workflows []domain.Workflow

	db := r.db.WithContext(ctx).Model(&domain.Workflow{}).
		Where("organization_id = ?", organizationID)

	if q.Search != "" {
		db = db.Where("name ILIKE ?", "%"+q.Search+"%")
	}

	db, total, err := Paginate(db, q)
	if err != nil {
		return nil, 0, err
	}

	err = db.Find(&workflows).Error
	if err != nil {
		return nil, 0, err
	}

	return workflows, total, nil
}

func (r *WorkflowRepository) FindByOrganizationIDAndWorkflowID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (domain.Workflow, error) {
	var workflow domain.Workflow
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", organizationID, workflowID).
		First(&workflow).Error
	if err != nil {
		return domain.Workflow{}, err
	}
	return workflow, nil
}
