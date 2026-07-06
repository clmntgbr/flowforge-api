package variable

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type VariableInUseError struct {
	Steps []entity.Step
}

func (e *VariableInUseError) Error() string {
	return "variable is used in workflow steps"
}

type DeleteVariableUseCase struct {
	variableRepo *repository.VariableRepository
	stepRepo     *repository.StepRepository
}

func NewDeleteVariableUseCase(variableRepo *repository.VariableRepository, stepRepo *repository.StepRepository) *DeleteVariableUseCase {
	return &DeleteVariableUseCase{
		variableRepo: variableRepo,
		stepRepo:     stepRepo,
	}
}

func (u *DeleteVariableUseCase) Execute(ctx context.Context, workflowID uuid.UUID, variableID uuid.UUID) error {
	variable, err := (*u.variableRepo).GetVariableByIDAndWorkflowID(ctx, workflowID, variableID)
	if err != nil {
		return err
	}

	steps, err := (*u.stepRepo).GetByWorkflowID(ctx, workflowID)
	if err != nil {
		return err
	}

	usedInSteps := findStepsUsingVariableKey(steps, variable.Key)
	if len(usedInSteps) > 0 {
		return &VariableInUseError{Steps: usedInSteps}
	}

	return (*u.variableRepo).Delete(ctx, variableID)
}
