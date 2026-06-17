package workflow_run

import (
	"context"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type ComputeSkippedStepsUseCase struct {
	stepRepo *repository.StepRepository
}

func NewComputeSkippedStepsUseCase(stepRepo *repository.StepRepository) *ComputeSkippedStepsUseCase {
	return &ComputeSkippedStepsUseCase{stepRepo: stepRepo}
}

func (u *ComputeSkippedStepsUseCase) Execute(ctx context.Context, workflowID uuid.UUID, executedStepIDs []string) ([]string, error) {
	steps, err := (*u.stepRepo).GetByWorkflowID(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	executed := make(map[string]struct{}, len(executedStepIDs))
	for _, stepID := range executedStepIDs {
		executed[stepID] = struct{}{}
	}

	skipped := make([]string, 0)
	for _, step := range steps {
		stepID := step.ID.String()
		if _, ok := executed[stepID]; !ok {
			skipped = append(skipped, stepID)
		}
	}

	return skipped, nil
}
