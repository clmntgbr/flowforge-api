package endpoint

import (
	"context"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type EndpointHasStepUseCase struct {
	stepRepo *repository.StepRepository
}

func NewEndpointHasStepUseCase(stepRepo *repository.StepRepository) *EndpointHasStepUseCase {
	return &EndpointHasStepUseCase{stepRepo: stepRepo}
}

func (u *EndpointHasStepUseCase) Execute(ctx context.Context, endpointID uuid.UUID) (bool, error) {
	hasSteps, err := (*u.stepRepo).HasStepsByEndpointID(ctx, endpointID)
	if err != nil {
		return false, err
	}
	return hasSteps, nil
}
