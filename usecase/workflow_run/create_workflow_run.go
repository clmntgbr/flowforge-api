package workflow_run

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"fmt"

	"github.com/google/uuid"
)

type CreateWorkflowRunUseCase struct {
	workflowRunRepo repository.WorkflowRunRepository
}

func NewCreateWorkflowRunUseCase(
	workflowRunRepo repository.WorkflowRunRepository,
) *CreateWorkflowRunUseCase {
	return &CreateWorkflowRunUseCase{
		workflowRunRepo: workflowRunRepo,
	}
}

func (u *CreateWorkflowRunUseCase) Execute(ctx context.Context, workflowID uuid.UUID) (*entity.WorkflowRun, error) {
	fmt.Println("🔄 Creating workflow run", workflowID)
	workflowRun := &entity.WorkflowRun{
		WorkflowID: workflowID,
		Status:     enum.WorkflowRunStatusPending,
	}

	err := u.workflowRunRepo.Create(ctx, workflowRun)
	if err != nil {
		return nil, err
	}

	return workflowRun, nil
}
