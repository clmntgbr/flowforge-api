package repository

import (
	"context"
	"forgeflow-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StepRepository struct {
	db *gorm.DB
}

func NewStepRepository(db *gorm.DB) *StepRepository {
	return &StepRepository{db: db}
}

func (r *StepRepository) Create(step *domain.Step) error {
	return r.db.Create(step).Error
}

func (r *StepRepository) Update(step *domain.Step) error {
	return r.db.Save(step).Error
}

func (r *StepRepository) Delete(step *domain.Step) error {
	return r.db.Delete(step).Error
}

func (r *StepRepository) FindByWorkflowID(ctx context.Context, workflowID uuid.UUID) ([]domain.Step, error) {
	var steps []domain.Step
	err := r.db.WithContext(ctx).
		Where("workflow_id = ?", workflowID).
		Find(&steps).Error
	return steps, err
}

func (r *StepRepository) DeleteByIDs(ctx context.Context, stepIDs []uuid.UUID) error {
	if len(stepIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Where("id IN ?", stepIDs).Delete(&domain.Step{}).Error
	})
}

func (r *StepRepository) FindByID(ctx context.Context, stepID uuid.UUID) (*domain.Step, error) {
	var step domain.Step
	err := r.db.WithContext(ctx).Where("id = ?", stepID).First(&step).Error
	if err != nil {
		return nil, err
	}
	return &step, nil
}

func (r *StepRepository) UpdatePositionAndIndex(ctx context.Context, stepID uuid.UUID, workflowID uuid.UUID, position domain.Position, index string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&domain.Step{}).
			Where("id = ? AND workflow_id = ?", stepID, workflowID).
			Updates(map[string]interface{}{
				"position_x": position.X,
				"position_y": position.Y,
				"index":      index,
			}).Error
	})
}

func (r *StepRepository) FindByOrganizationIDAndStepID(ctx context.Context, organizationID uuid.UUID, stepID uuid.UUID) (*domain.Step, error) {
	var step domain.Step
	err := r.db.WithContext(ctx).
		Joins("JOIN workflows ON steps.workflow_id = workflows.id").
		Where("workflows.organization_id = ? AND steps.id = ?", organizationID, stepID).
		First(&step).Error
	if err != nil {
		return nil, err
	}
	return &step, nil
}
