package step

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	stepDTO "flowforge-api/infrastructure/step"

	"github.com/google/uuid"
)

type UpdateStepUseCase struct {
	stepRepo *repository.StepRepository
}

func NewUpdateStepUseCase(stepRepo *repository.StepRepository) *UpdateStepUseCase {
	return &UpdateStepUseCase{stepRepo: stepRepo}
}

func (u *UpdateStepUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, id uuid.UUID, input stepDTO.UpdateStepInput) (entity.Step, error) {
	step, err := (*u.stepRepo).GetByIDAndOrganizationIDAndWorkflowID(ctx, organizationID, workflowID, id)
	if err != nil {
		return entity.Step{}, err
	}

	step.Name = input.Name
	step.Description = input.Description
	step.Timeout = input.Timeout
	step.Query = input.Query
	step.Header = input.Header
	step.Body = input.Body
	step.RetryOnFailure = input.RetryOnFailure
	step.RetryCount = input.RetryCount
	step.RetryDelay = input.RetryDelay

	if err := (*u.stepRepo).Update(ctx, &step); err != nil {
		return entity.Step{}, err
	}

	return step, nil
}
