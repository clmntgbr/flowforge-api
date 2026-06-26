package variable

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type GetVariablesByWorkflowIDUseCase struct {
	variableRepo *repository.VariableRepository
}

func NewGetVariablesByWorkflowIDUseCase(variableRepo *repository.VariableRepository) *GetVariablesByWorkflowIDUseCase {
	return &GetVariablesByWorkflowIDUseCase{variableRepo: variableRepo}
}

func (u *GetVariablesByWorkflowIDUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) ([]entity.Variable, error) {
	variables, err := (*u.variableRepo).GetVariablesByWorkflowID(ctx, organizationID, workflowID)
	if err != nil {
		return []entity.Variable{}, err
	}

	return variables, nil
}
