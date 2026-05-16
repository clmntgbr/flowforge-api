package repository

import (
	"context"
	"flowforge-api/domain/entity"

	"github.com/google/uuid"
)

type InsightRepository interface {
	Create(ctx context.Context, insight *entity.Insight) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, insight *entity.Insight) error
}
