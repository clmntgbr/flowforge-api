package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type GetEndpointUseCase struct {
	endpointRepo repository.EndpointRepository
}

func NewGetEndpointUseCase(endpointRepo repository.EndpointRepository) *GetEndpointUseCase {
	return &GetEndpointUseCase{endpointRepo: endpointRepo}
}

func (u *GetEndpointUseCase) Execute(ctx context.Context, organizationID uuid.UUID, endpointID uuid.UUID) (entity.Endpoint, error) {
	endpoint, err := u.endpointRepo.GetByIDAndOrganizationID(ctx, endpointID, organizationID)
	if err != nil {
		return entity.Endpoint{}, err
	}

	return endpoint, nil
}
