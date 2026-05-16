package repository

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
)

type WorkflowRunRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.WorkflowRun, error)
	GetByWorkflowID(ctx context.Context, workflowID uuid.UUID, query paginate.PaginateQuery) ([]entity.WorkflowRun, int64, error)
	GetByWorkflowIDAndNotEnded(ctx context.Context, workflowID uuid.UUID) (*entity.WorkflowRun, error)
	Create(ctx context.Context, workflowRun *entity.WorkflowRun) error
	Update(ctx context.Context, workflowRun *entity.WorkflowRun) error
	Delete(ctx context.Context, id uuid.UUID) error
}
