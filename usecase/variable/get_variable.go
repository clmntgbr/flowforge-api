package variable

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type GetVariableByIDUseCase struct {
	variableRepo *repository.VariableRepository
}

func NewGetVariableByIDUseCase(variableRepo *repository.VariableRepository) *GetVariableByIDUseCase {
	return &GetVariableByIDUseCase{variableRepo: variableRepo}
}

func (u *GetVariableByIDUseCase) Execute(ctx context.Context, workflowID uuid.UUID, variableID uuid.UUID) (entity.Variable, error) {
	variable, err := (*u.variableRepo).GetVariableByIDAndWorkflowID(ctx, workflowID, variableID)
	if err != nil {
		return entity.Variable{}, err
	}

	return variable, nil
}
