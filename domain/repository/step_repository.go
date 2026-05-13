package repository

import (
	"context"
	"flowforge-api/domain/entity"

	"github.com/google/uuid"
)

type StepRepository interface {
	Create(ctx context.Context, step *entity.Step) error
	Update(ctx context.Context, step *entity.Step) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByWorkflowID(ctx context.Context, workflowID uuid.UUID) ([]entity.Step, error)
	GetByIDAndOrganizationIDAndWorkflowID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, id uuid.UUID) (entity.Step, error)
	DeleteByIDs(ctx context.Context, ids []uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Step, error)
	UpdatePositionAndIndex(ctx context.Context, id uuid.UUID, workflowID uuid.UUID, position entity.Position, index string, executionOrder int) error
}
