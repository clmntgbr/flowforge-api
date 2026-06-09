package workflow_run

import (
	"context"
	"errors"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type IsCanceledWorkflowRunUseCase struct {
	workflowRunRepo *repository.WorkflowRunRepository
}

func NewIsCanceledWorkflowRunUseCase(
	workflowRunRepo *repository.WorkflowRunRepository,
) *IsCanceledWorkflowRunUseCase {
	return &IsCanceledWorkflowRunUseCase{
		workflowRunRepo: workflowRunRepo,
	}
}

func (u *IsCanceledWorkflowRunUseCase) Execute(ctx context.Context, workflowRunID string) error {
	workflowRunUUID, err := uuid.Parse(workflowRunID)
	if err != nil {
		return err
	}

	workflowRun, err := (*u.workflowRunRepo).GetByID(ctx, workflowRunUUID)
	if err != nil {
		return err
	}

	if workflowRun.Status == enum.WorkflowRunStatusCanceled {
		return errors.New("workflow run is canceled")
	}

	return nil
}
