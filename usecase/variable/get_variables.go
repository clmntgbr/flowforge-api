package variable

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type GetVariablesByWorkflowIDUseCase struct {
	variableRepo *repository.VariableRepository
}

func NewGetVariablesByWorkflowIDUseCase(variableRepo *repository.VariableRepository) *GetVariablesByWorkflowIDUseCase {
	return &GetVariablesByWorkflowIDUseCase{variableRepo: variableRepo}
}

func (u *GetVariablesByWorkflowIDUseCase) Execute(ctx context.Context, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.Variable, int64, error) {
	variables, total, err := (*u.variableRepo).ListByWorkflowID(ctx, workflowID, query)
	if err != nil {
		return []entity.Variable{}, 0, err
	}

	return variables, total, nil
}
