package gorm

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type variableRepository struct {
	db *gorm.DB
}

func NewVariableRepository(db *gorm.DB) repository.VariableRepository {
	return &variableRepository{db: db}
}

func (r *variableRepository) Create(ctx context.Context, variable *entity.Variable) error {
	return dbWithContext(ctx, r.db).Create(variable).Error
}

func (r *variableRepository) Update(ctx context.Context, variable *entity.Variable) error {
	return dbWithContext(ctx, r.db).Save(variable).Error
}

func (r *variableRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Variable{}, id).Error
}

func (r *variableRepository) GetVariablesByWorkflowID(ctx context.Context, workflowID uuid.UUID) ([]entity.Variable, error) {
	var variables []entity.Variable

	db := dbWithContext(ctx, r.db).Model(&entity.Variable{}).
		Where("variables.workflow_id = ?", workflowID).
		Preload("Step").
		Preload("Step.Endpoint")

	err := db.Find(&variables).Error
	if err != nil {
		return nil, err
	}

	return variables, nil
}

func (r *variableRepository) ListByWorkflowID(ctx context.Context, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.Variable, int64, error) {
	var variables []entity.Variable

	db := dbWithContext(ctx, r.db).Model(&entity.Variable{}).
		Where("variables.workflow_id = ?", workflowID).
		Preload("Step").
		Preload("Step.Endpoint")

	if query.Search != "" {
		db = db.Where("variables.name ILIKE ? OR variables.key ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}

	db, total, err := Paginate(db, query)
	if err != nil {
		return nil, 0, err
	}

	err = db.Find(&variables).Error
	if err != nil {
		return nil, 0, err
	}

	return variables, total, nil
}

func (r *variableRepository) GetVariableByIDAndWorkflowID(ctx context.Context, workflowID uuid.UUID, variableID uuid.UUID) (entity.Variable, error) {
	var variable entity.Variable

	db := dbWithContext(ctx, r.db).Model(&entity.Variable{}).
		Where("variables.id = ?", variableID).
		Where("variables.workflow_id = ?", workflowID).
		Preload("Step").
		Preload("Step.Endpoint")

	err := db.Find(&variable).Error
	if err != nil {
		return entity.Variable{}, err
	}

	return variable, nil
}

func (r *variableRepository) GetVariableByWorkflowIDAndKey(ctx context.Context, workflowID uuid.UUID, key string) (*entity.Variable, error) {
	var variable entity.Variable

	err := dbWithContext(ctx, r.db).
		Where("workflow_id = ?", workflowID).
		Where("key = ?", key).
		First(&variable).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &variable, nil
}

func (r *variableRepository) GetVariablesByStepID(ctx context.Context, stepID uuid.UUID) ([]entity.Variable, error) {
	var variables []entity.Variable

	err := dbWithContext(ctx, r.db).Model(&entity.Variable{}).
		Where("variables.step_id = ?", stepID).
		Preload("Step").
		Preload("Step.Endpoint").
		Find(&variables).Error
	if err != nil {
		return nil, err
	}

	return variables, nil
}
