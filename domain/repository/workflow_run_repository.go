package repository

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type WorkflowRunRepository interface {
	GetByWorkflowID(ctx context.Context, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.WorkflowRun, int64, error)
}
