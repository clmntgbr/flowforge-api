package workflow_run

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type CreateWorkflowRunUseCase struct {
	workflowRunRepo *repository.WorkflowRunRepository
}

func NewCreateWorkflowRunUseCase(
	workflowRunRepo *repository.WorkflowRunRepository,
) *CreateWorkflowRunUseCase {
	return &CreateWorkflowRunUseCase{
		workflowRunRepo: workflowRunRepo,
	}
}

func (u *CreateWorkflowRunUseCase) Execute(ctx context.Context, workflowID uuid.UUID, totalSteps int) (*entity.WorkflowRun, error) {
	workflowRun := &entity.WorkflowRun{
		WorkflowID: workflowID,
		Status:     enum.WorkflowRunStatusPending,
		Statuses:   []enum.WorkflowRunStatus{enum.WorkflowRunStatusPending},
		TotalSteps: totalSteps,
	}

	err := (*u.workflowRunRepo).Create(ctx, workflowRun)
	if err != nil {
		return nil, err
	}

	return workflowRun, nil
}
