package variable

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	variableDTO "flowforge-api/infrastructure/variable"

	"github.com/google/uuid"
)

type CreateVariableUseCase struct {
	variableRepo *repository.VariableRepository
	workflowRepo *repository.WorkflowRepository
}

func NewCreateVariableUseCase(variableRepo *repository.VariableRepository, workflowRepo *repository.WorkflowRepository) *CreateVariableUseCase {
	return &CreateVariableUseCase{variableRepo: variableRepo, workflowRepo: workflowRepo}
}

func (u *CreateVariableUseCase) Execute(ctx context.Context, workflowID uuid.UUID, organizationID uuid.UUID, request variableDTO.CreateVariableInput) (entity.Variable, error) {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, workflowID, organizationID)
	if err != nil {
		return entity.Variable{}, err
	}

	variable := &entity.Variable{
		Name:        request.Name,
		Description: request.Description,
		Path:        request.Path,
		StepID:      request.StepID,
		WorkflowID:  workflow.ID,
	}

	err = (*u.variableRepo).Create(ctx, variable)
	if err != nil {
		return entity.Variable{}, err
	}

	return *variable, nil
}
