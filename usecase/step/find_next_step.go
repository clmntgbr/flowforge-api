package step

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type FindNextStepUseCase struct {
	stepRepo *repository.StepRepository
}

func NewFindNextStepUseCase(stepRepo *repository.StepRepository) *FindNextStepUseCase {
	return &FindNextStepUseCase{stepRepo: stepRepo}
}

func (u *FindNextStepUseCase) Execute(ctx context.Context, workflowID uuid.UUID, executedStepIDs []string) (*entity.Step, error) {
	return (*u.stepRepo).GetNextStepByWorkflowID(ctx, workflowID, 0, executedStepIDs)
}
