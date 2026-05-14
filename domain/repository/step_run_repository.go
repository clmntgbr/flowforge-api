package repository

import (
	"context"
	"flowforge-api/domain/entity"

	"github.com/google/uuid"
)

type StepRunRepository interface {
	Create(ctx context.Context, stepRun *entity.StepRun) error
	Update(ctx context.Context, stepRun *entity.StepRun) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.StepRun, error)
	GetByWorkflowRunID(ctx context.Context, workflowRunID uuid.UUID) (*entity.StepRun, error)
}
