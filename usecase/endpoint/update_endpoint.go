package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	endpointDTO "flowforge-api/infrastructure/endpoint"

	"github.com/google/uuid"
)

type UpdateEndpointUseCase struct {
	endpointRepo *repository.EndpointRepository
}

func NewUpdateEndpointUseCase(endpointRepo *repository.EndpointRepository) *UpdateEndpointUseCase {
	return &UpdateEndpointUseCase{endpointRepo: endpointRepo}
}

func (u *UpdateEndpointUseCase) Execute(ctx context.Context, organizationID uuid.UUID, endpointID uuid.UUID, input endpointDTO.UpdateEndpointInput) (entity.Endpoint, error) {
	endpoint, err := (*u.endpointRepo).GetByIDAndOrganizationID(ctx, endpointID, organizationID)
	if err != nil {
		return entity.Endpoint{}, err
	}

	endpoint.Name = input.Name
	endpoint.BaseURI = input.BaseURI
	endpoint.Path = input.Path
	endpoint.Method = input.Method
	endpoint.Timeout = input.Timeout
	endpoint.Query = input.Query
	endpoint.Header = input.Header
	endpoint.Body = input.Body
	endpoint.RetryOnFailure = input.RetryOnFailure
	endpoint.RetryCount = input.RetryCount
	endpoint.RetryDelay = input.RetryDelay

	if err := (*u.endpointRepo).Update(ctx, &endpoint); err != nil {
		return entity.Endpoint{}, err
	}

	return endpoint, nil
}
