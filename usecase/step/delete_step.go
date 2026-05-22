package step

import (
	"context"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type DeleteStepUseCase struct {
	stepRepo      *repository.StepRepository
	connexionRepo *repository.ConnexionRepository
	workflowRepo  *repository.WorkflowRepository
}

func NewDeleteStepUseCase(
	stepRepo *repository.StepRepository,
	connexionRepo *repository.ConnexionRepository,
	workflowRepo *repository.WorkflowRepository,
) *DeleteStepUseCase {
	return &DeleteStepUseCase{
		stepRepo:      stepRepo,
		connexionRepo: connexionRepo,
		workflowRepo:  workflowRepo,
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

	if err := (*u.connexionRepo).DeleteByStepID(ctx, step.ID); err != nil {
		return err
	}

	if err := (*u.stepRepo).Delete(ctx, step.ID); err != nil {
		return err
	}

	return nil
}
