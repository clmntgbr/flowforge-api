package repository

import (
	"context"
	"flowforge-api/domain/entity"
	endpointDTO "flowforge-api/infrastructure/endpoint"

	"github.com/google/uuid"
)

type EndpointRepository interface {
	List(ctx context.Context, organizationID uuid.UUID, query endpointDTO.PaginateEndpointQuery) ([]entity.Endpoint, int64, error)
	Create(ctx context.Context, endpoint *entity.Endpoint) error
	Update(ctx context.Context, endpoint *entity.Endpoint) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (entity.Endpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Endpoint, error)
}
