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
	GetByIDAndOrganizationIDAndWorkflowID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID, id uuid.UUID) (entity.Step, error)
}
