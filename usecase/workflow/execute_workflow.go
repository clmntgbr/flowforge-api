package workflow

import (
	"context"
	"flowforge-api/domain/repository"
)

type ExecuteWorkflowUseCase struct {
	workflowRepo repository.WorkflowRepository
}

func NewExecuteWorkflowUseCase(workflowRepo repository.WorkflowRepository) *ExecuteWorkflowUseCase {
	return &ExecuteWorkflowUseCase{workflowRepo: workflowRepo}
}

func (u *ExecuteWorkflowUseCase) Execute(ctx context.Context) error {
	return nil
}
