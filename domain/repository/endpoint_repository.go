package repository

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type EndpointRepository interface {
	List(ctx context.Context, organizationID uuid.UUID, query paginate.PaginateQuery) ([]entity.Endpoint, int64, error)
	Create(ctx context.Context, endpoint *entity.Endpoint) error
	Update(ctx context.Context, endpoint *entity.Endpoint) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (entity.Endpoint, error)
}
