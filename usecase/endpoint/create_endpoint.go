package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	endpointDTO "flowforge-api/infrastructure/endpoint"

	"github.com/google/uuid"
)

type CreateEndpointUseCase struct {
	endpointRepo *repository.EndpointRepository
}

func NewCreateEndpointUseCase(endpointRepo *repository.EndpointRepository) *CreateEndpointUseCase {
	return &CreateEndpointUseCase{endpointRepo: endpointRepo}
}

func (u *CreateEndpointUseCase) Execute(ctx context.Context, organizationID uuid.UUID, input endpointDTO.CreateEndpointInput) (entity.Endpoint, error) {
	endpoint := &entity.Endpoint{
		Name:           input.Name,
		OrganizationID: organizationID,
		BaseURI:        input.BaseURI,
		Path:           input.Path,
		Method:         input.Method,
		Timeout:        input.Timeout,
		Query:          input.Query,
		Header:         input.Header,
		Body:           input.Body,
		RetryOnFailure: input.RetryOnFailure,
		RetryCount:     input.RetryCount,
		RetryDelay:     input.RetryDelay,
	}

	if err := (*u.endpointRepo).Create(ctx, endpoint); err != nil {
		return entity.Endpoint{}, err
	}

	return *endpoint, nil
}
