package gorm

import (
	"context"
	"errors"
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
	return dbWithContext(ctx, r.db).Create(step).Error
}

func (r *stepRepository) Update(ctx context.Context, step *entity.Step) error {
	return dbWithContext(ctx, r.db).Save(step).Error
}

func (r *stepRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Step{}, id).Error
}

func (r *stepRepository) GetByIDAndOrganizationIDAndWorkflowID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, id uuid.UUID) (entity.Step, error) {
	var step entity.Step

	err := dbWithContext(ctx, r.db).
		Joins("JOIN workflows ON workflows.id = steps.workflow_id").
		Where("steps.id = ? AND workflows.organization_id = ? AND workflows.id = ?", id, organizationID, workflowID).
		First(&step).Error
	if err != nil {
		return entity.Step{}, err
	}
	return step, nil
}

func (r *stepRepository) DeleteByIDs(ctx context.Context, stepIDs []uuid.UUID) error {
	if len(stepIDs) == 0 {
		return nil
	}
	return dbWithContext(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		return tx.Where("id IN ?", stepIDs).Delete(&entity.Step{}).Error
	})
}

func (r *stepRepository) GetByWorkflowID(ctx context.Context, workflowID uuid.UUID) ([]entity.Step, error) {
	var steps []entity.Step
	err := dbWithContext(ctx, r.db).
		Where("workflow_id = ?", workflowID).
		Find(&steps).Error
	return steps, err
}

func (r *stepRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Step, error) {
	var step entity.Step
	err := dbWithContext(ctx, r.db).
		Where("id = ?", id).
		Preload("Endpoint").
		First(&step).Error

	if err != nil {
		return nil, err
	}

	return &step, nil
}

func (r *stepRepository) UpdatePositionAndIndex(ctx context.Context, stepID uuid.UUID, workflowID uuid.UUID, position entity.Position, index string, executionOrder int) error {
	return dbWithContext(ctx, r.db).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&entity.Step{}).
			Where("id = ? AND workflow_id = ?", stepID, workflowID).
			Updates(map[string]interface{}{
				"position_x":      position.X,
				"position_y":      position.Y,
				"index":           index,
				"execution_order": executionOrder,
			}).Error
	})
}

func (r *stepRepository) GetFirstStepByWorkflowID(ctx context.Context, workflowID uuid.UUID) (*entity.Step, error) {
	var step entity.Step
	err := dbWithContext(ctx, r.db).
		Where("workflow_id = ?", workflowID).
		Order("execution_order ASC, id ASC").
		Preload("Endpoint").
		First(&step).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &step, nil
}
