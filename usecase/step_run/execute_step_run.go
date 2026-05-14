package step_run

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"time"
)

type ExecuteStepRunUseCase struct {
	stepRunRepo repository.StepRunRepository
	stepRepo    repository.StepRepository
}

func NewExecuteStepRunUseCase(
	stepRunRepo repository.StepRunRepository,
	stepRepo repository.StepRepository,
) *ExecuteStepRunUseCase {
	return &ExecuteStepRunUseCase{
		stepRunRepo: stepRunRepo,
		stepRepo:    stepRepo,
	}
}

func (u *ExecuteStepRunUseCase) Execute(ctx context.Context, stepRun *entity.StepRun) (entity.StepRun, error) {
	stepRun.Status = enum.StepRunStatusRunning
	startedAt := time.Now().UTC()
	stepRun.StartedAt = &startedAt

	err := u.stepRunRepo.Update(ctx, stepRun)
	if err != nil {
		return entity.StepRun{}, err
	}

	return *stepRun, nil
}
