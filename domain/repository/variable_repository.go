package repository

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type VariableRepository interface {
	Create(ctx context.Context, variable *entity.Variable) error
	Update(ctx context.Context, variable *entity.Variable) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetVariablesByWorkflowID(ctx context.Context, workflowID uuid.UUID) ([]entity.Variable, error)
	ListByWorkflowID(ctx context.Context, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.Variable, int64, error)
	GetVariableByIDAndWorkflowID(ctx context.Context, workflowID uuid.UUID, variableID uuid.UUID) (entity.Variable, error)
	GetVariableByWorkflowIDAndKey(ctx context.Context, workflowID uuid.UUID, key string) (*entity.Variable, error)
}
