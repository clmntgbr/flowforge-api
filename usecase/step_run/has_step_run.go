package step_run

import (
	"context"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type HasStepRunUseCase struct {
	stepRunRepo *repository.StepRunRepository
}

func NewHasStepRunUseCase(
	stepRunRepo *repository.StepRunRepository,
) *HasStepRunUseCase {
	return &HasStepRunUseCase{
		stepRunRepo: stepRunRepo,
	}
}

func (u *HasStepRunUseCase) Execute(ctx context.Context, workflowRunID uuid.UUID) bool {
	stepRun, err := (*u.stepRunRepo).GetByWorkflowRunID(ctx, workflowRunID)
	if err != nil {
		return false
	}

	if stepRun == nil {
		return false
	}

	return true
}
