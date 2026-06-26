package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

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

func (r *variableRepository) GetVariablesByWorkflowID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) ([]entity.Variable, error) {
	var variables []entity.Variable

	db := dbWithContext(ctx, r.db).Model(&entity.Variable{}).
		Joins("JOIN workflows ON workflows.id = variables.workflow_id").
		Where("workflows.organization_id = ? AND variables.workflow_id = ?", organizationID, workflowID)

	err := db.Find(&variables).Error
	if err != nil {
		return nil, err
	}

	return variables, nil
}
