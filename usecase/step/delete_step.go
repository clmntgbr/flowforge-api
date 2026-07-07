package step

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type StepInUseError struct {
	Variables []entity.Variable
}

func (e *StepInUseError) Error() string {
	return "step is used in workflow variables"
}

type DeleteStepUseCase struct {
	stepRepo       *repository.StepRepository
	connexionRepo  *repository.ConnexionRepository
	workflowRepo   *repository.WorkflowRepository
	variableRepo   *repository.VariableRepository
}

func NewDeleteStepUseCase(
	stepRepo *repository.StepRepository,
	connexionRepo *repository.ConnexionRepository,
	workflowRepo *repository.WorkflowRepository,
	variableRepo *repository.VariableRepository,
) *DeleteStepUseCase {
	return &DeleteStepUseCase{
		stepRepo:      stepRepo,
		connexionRepo: connexionRepo,
		workflowRepo:  workflowRepo,
		variableRepo:  variableRepo,
	}
}

func (u *DeleteStepUseCase) Execute(
	ctx context.Context,
	organizationID uuid.UUID,
	workflowID uuid.UUID,
	stepUUID uuid.UUID,
) error {
	step, err := (*u.stepRepo).GetByIDAndOrganizationIDAndWorkflowID(ctx, organizationID, workflowID, stepUUID)
	if err != nil {
		return err
	}

	if step.IsEnabled {
		variables, err := (*u.variableRepo).GetVariablesByStepID(ctx, step.ID)
		if err != nil {
			return err
		}

		if len(variables) > 0 {
			return &StepInUseError{Variables: variables}
		}
	}

	if err := (*u.connexionRepo).DeleteByStepID(ctx, step.ID); err != nil {
		return err
	}

	step.IsEnabled = false
	if err := (*u.stepRepo).Update(ctx, &step); err != nil {
		return err
	}

	return nil
}
