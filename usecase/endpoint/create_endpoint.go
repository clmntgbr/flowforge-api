package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	endpointDTO "flowforge-api/infrastructure/endpoint"
	"log"

	"github.com/google/uuid"
)

type CreateEndpointUseCase struct {
	endpointRepo repository.EndpointRepository
}

func NewCreateEndpointUseCase(endpointRepo repository.EndpointRepository) *CreateEndpointUseCase {
	return &CreateEndpointUseCase{endpointRepo: endpointRepo}
}

func (u *CreateEndpointUseCase) Execute(ctx context.Context, organizationID uuid.UUID, input endpointDTO.CreateEndpointInput) (entity.Endpoint, error) {

	log.Println("input: ", input)
	log.Println("organizationID: ", organizationID)
	log.Println("Body: ", input.Body)
	log.Println("Query: ", input.Query)
	log.Println("Header: ", input.Header)

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

	if err := u.endpointRepo.Create(ctx, endpoint); err != nil {
		log.Println("error creating endpoint: ", err)
		return entity.Endpoint{}, err
	}

	return *endpoint, nil
}
