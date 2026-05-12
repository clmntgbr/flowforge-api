package repository

import (
	"context"
	"flowforge-api/domain/entity"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Organization, error)
	Create(ctx context.Context, organization *entity.Organization) error
	Update(ctx context.Context, organization *entity.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
}
