package step

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type GetStepUseCase struct {
	stepRepo repository.StepRepository
}

func NewGetStepUseCase(stepRepo repository.StepRepository) *GetStepUseCase {
	return &GetStepUseCase{stepRepo: stepRepo}
}

func (u *GetStepUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, id uuid.UUID) (entity.Step, error) {
	step, err := u.stepRepo.GetByIDAndOrganizationIDAndWorkflowID(ctx, organizationID, workflowID, id)
	if err != nil {
		return entity.Step{}, err
	}

	return step, nil
}
