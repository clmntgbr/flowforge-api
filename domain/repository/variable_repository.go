package repository

import (
	"context"
	"flowforge-api/domain/entity"

	"github.com/google/uuid"
)

type VariableRepository interface {
	Create(ctx context.Context, variable *entity.Variable) error
	Update(ctx context.Context, variable *entity.Variable) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetVariablesByWorkflowID(ctx context.Context, workflowID uuid.UUID) ([]entity.Variable, error)
	GetVariableByIDAndWorkflowID(ctx context.Context, workflowID uuid.UUID, variableID uuid.UUID) (entity.Variable, error)
	GetVariableByWorkflowIDAndKey(ctx context.Context, workflowID uuid.UUID, key string) (*entity.Variable, error)
}
