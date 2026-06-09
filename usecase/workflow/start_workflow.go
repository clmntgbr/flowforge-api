package workflow

import (
	"context"
	"errors"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type StartWorkflowUseCase struct {
	workflowRepo       *repository.WorkflowRepository
	workflowRunRepo    *repository.WorkflowRunRepository
	runWorkflowUseCase *RunWorkflowUseCase
}

func NewStartWorkflowUseCase(
	workflowRepo *repository.WorkflowRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
	runWorkflowUseCase *RunWorkflowUseCase,
) *StartWorkflowUseCase {
	return &StartWorkflowUseCase{
		workflowRepo:       workflowRepo,
		workflowRunRepo:    workflowRunRepo,
		runWorkflowUseCase: runWorkflowUseCase,
	}
}

func (u *StartWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) error {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return err
	}

	if workflow.Status == enum.WorkflowStatusInactive {
		return errors.New("workflow is inactive")
	}

	err = u.runWorkflowUseCase.Execute(ctx, workflow)
	if err != nil {
		return err
	}

	return nil
}
