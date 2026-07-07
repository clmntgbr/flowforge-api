package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type EndpointInUseError struct {
	Steps []entity.Step
}

func (e *EndpointInUseError) Error() string {
	return "endpoint is used in workflow steps"
}

type DeleteEndpointUseCase struct {
	endpointRepo *repository.EndpointRepository
	stepRepo     *repository.StepRepository
}

func NewDeleteEndpointUseCase(
	endpointRepo *repository.EndpointRepository,
	stepRepo *repository.StepRepository,
) *DeleteEndpointUseCase {
	return &DeleteEndpointUseCase{
		endpointRepo: endpointRepo,
		stepRepo:     stepRepo,
	}
}

func (u *DeleteEndpointUseCase) Execute(ctx context.Context, endpointID uuid.UUID) error {
	steps, err := (*u.stepRepo).GetEnabledStepsByEndpointID(ctx, endpointID)
	if err != nil {
		return err
	}

	if len(steps) > 0 {
		return &EndpointInUseError{Steps: steps}
	}

	return (*u.endpointRepo).Delete(ctx, endpointID)
}
