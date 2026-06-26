package variable

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	variableDTO "flowforge-api/infrastructure/variable"

	"github.com/google/uuid"
)

type UpdateVariableUseCase struct {
	variableRepo *repository.VariableRepository
}

func NewUpdateVariableUseCase(variableRepo *repository.VariableRepository) *UpdateVariableUseCase {
	return &UpdateVariableUseCase{variableRepo: variableRepo}
}

func (u *UpdateVariableUseCase) Execute(ctx context.Context, workflowID uuid.UUID, variableID uuid.UUID, request variableDTO.UpdateVariableInput) (entity.Variable, error) {

	variable, err := (*u.variableRepo).GetVariableByIDAndWorkflowID(ctx, workflowID, variableID)
	if err != nil {
		return entity.Variable{}, err
	}

	variable.Name = request.Name
	variable.Path = request.Path
	variable.Description = request.Description
	variable.StepID = request.StepID

	err = (*u.variableRepo).Update(ctx, &variable)
	if err != nil {
		return entity.Variable{}, err
	}

	return variable, nil
}
