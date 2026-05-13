package repository

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkflowRepository interface {
	List(ctx context.Context, organizationID uuid.UUID, query paginate.PaginateQuery) ([]entity.Workflow, int64, error)
	Create(ctx context.Context, workflow *entity.Workflow) error
	Update(ctx context.Context, workflow *entity.Workflow) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByIDAndOrganizationID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (entity.Workflow, error)
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}
