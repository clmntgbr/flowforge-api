package endpoint

import (
	"context"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type DeleteEndpointUseCase struct {
	endpointRepo *repository.EndpointRepository
}

func NewDeleteEndpointUseCase(endpointRepo *repository.EndpointRepository) *DeleteEndpointUseCase {
	return &DeleteEndpointUseCase{endpointRepo: endpointRepo}
}

func (u *DeleteEndpointUseCase) Execute(ctx context.Context, endpointID uuid.UUID) error {
	if err := (*u.endpointRepo).Delete(ctx, endpointID); err != nil {
		return err
	}
	return nil
}
