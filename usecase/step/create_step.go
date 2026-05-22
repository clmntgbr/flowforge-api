package step

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/domain/types"

	"github.com/google/uuid"
)

type CreateStepUseCase struct {
	stepRepo *repository.StepRepository
}

func NewCreateStepUseCase(stepRepo *repository.StepRepository) *CreateStepUseCase {
	return &CreateStepUseCase{stepRepo: stepRepo}
}

func (u *CreateStepUseCase) Execute(
	ctx context.Context,
	workflowID uuid.UUID,
	stepUUID uuid.UUID,
	endpoint entity.Endpoint,
	position entity.Position,
	index string,
	executionOrder int,
	endpointUUID uuid.UUID,
) (entity.Step, error) {

	if len(endpoint.Header) == 0 {
		endpoint.Header = types.Header{
			{
				ID:    uuid.New().String(),
				Key:   "Content-Type",
				Value: "application/json",
			},
		}
	}

	step := &entity.Step{
		ID:             stepUUID,
		Name:           endpoint.Name,
		Description:    endpoint.Description,
		Timeout:        endpoint.Timeout,
		Query:          endpoint.Query,
		Header:         endpoint.Header,
		Body:           endpoint.Body,
		Position:       position,
		Index:          index,
		ExecutionOrder: executionOrder,
		EndpointID:     endpointUUID,
		WorkflowID:     workflowID,
		RetryOnFailure: endpoint.RetryOnFailure,
		RetryCount:     endpoint.RetryCount,
		RetryDelay:     endpoint.RetryDelay,
		IsEnabled:      true,
	}

	if err := (*u.stepRepo).Create(ctx, step); err != nil {
		return entity.Step{}, err
	}

	return *step, nil
}
