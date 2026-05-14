package step_run

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type CreateStepRunUseCase struct {
	stepRunRepo repository.StepRunRepository
	stepRepo    repository.StepRepository
}

func NewCreateStepRunUseCase(
	stepRunRepo repository.StepRunRepository,
	stepRepo repository.StepRepository,
) *CreateStepRunUseCase {
	return &CreateStepRunUseCase{
		stepRunRepo: stepRunRepo,
		stepRepo:    stepRepo,
	}
}

func (u *CreateStepRunUseCase) Execute(ctx context.Context, workflowRunID uuid.UUID, stepID uuid.UUID) (entity.StepRun, error) {
	stepRun := &entity.StepRun{
		StepID:        stepID,
		WorkflowRunID: workflowRunID,
		Status:        enum.StepRunStatusPending,
	}

	err := u.stepRunRepo.Create(ctx, stepRun)
	if err != nil {
		return entity.StepRun{}, err
	}

	step, err := u.stepRepo.GetByID(ctx, stepRun.StepID)
	if err != nil {
		return entity.StepRun{}, err
	}

	stepRun.Step = *step
	return *stepRun, nil
}
