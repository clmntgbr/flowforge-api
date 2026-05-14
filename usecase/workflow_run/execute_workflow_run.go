package workflow_run

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"time"
)

type ExecuteWorkflowRunUseCase struct {
	workflowRunRepo repository.WorkflowRunRepository
}

func NewExecuteWorkflowRunUseCase(
	workflowRunRepo repository.WorkflowRunRepository,
) *ExecuteWorkflowRunUseCase {
	return &ExecuteWorkflowRunUseCase{
		workflowRunRepo: workflowRunRepo,
	}
}

func (u *ExecuteWorkflowRunUseCase) Execute(ctx context.Context, workflowRun entity.WorkflowRun) (entity.WorkflowRun, error) {
	workflowRun.Status = enum.WorkflowRunStatusRunning
	startedAt := time.Now().UTC()
	workflowRun.StartedAt = &startedAt

	err := u.workflowRunRepo.Update(ctx, &workflowRun)
	if err != nil {
		return entity.WorkflowRun{}, err
	}

	return workflowRun, nil
}
