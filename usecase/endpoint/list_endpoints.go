package endpoint

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	endpointDTO "flowforge-api/infrastructure/endpoint"

	"github.com/google/uuid"
)

type ListEndpointsUseCase struct {
	endpointRepo *repository.EndpointRepository
}

func NewListEndpointsUseCase(endpointRepo *repository.EndpointRepository) *ListEndpointsUseCase {
	return &ListEndpointsUseCase{endpointRepo: endpointRepo}
}

func (u *ListEndpointsUseCase) Execute(ctx context.Context, organizationID uuid.UUID, query endpointDTO.PaginateEndpointQuery) ([]entity.Endpoint, int64, error) {
	endpoints, total, err := (*u.endpointRepo).List(ctx, organizationID, query)
	if err != nil {
		return []entity.Endpoint{}, 0, err
	}

	return endpoints, total, nil
}
