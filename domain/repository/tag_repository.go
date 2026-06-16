package repository

import (
	"context"
	"flowforge-api/domain/entity"

	"github.com/google/uuid"
)

type TagRepository interface {
	Get(ctx context.Context, organizationID uuid.UUID, tagID uuid.UUID) (entity.Tag, error)
	List(ctx context.Context, organizationID uuid.UUID) ([]entity.Tag, error)
	Create(ctx context.Context, tag *entity.Tag) error
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) error
}
