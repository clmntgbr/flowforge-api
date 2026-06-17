package workflow

import (
	"context"
	"errors"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/mercure"
	"fmt"

	"github.com/google/uuid"
)

type StartWorkflowUseCase struct {
	workflowRepo       *repository.WorkflowRepository
	workflowRunRepo    *repository.WorkflowRunRepository
	runWorkflowUseCase *RunWorkflowUseCase
	mercurePublisher   *mercure.Publisher
}

func NewStartWorkflowUseCase(
	workflowRepo *repository.WorkflowRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
	runWorkflowUseCase *RunWorkflowUseCase,
	mercurePublisher *mercure.Publisher,
) *StartWorkflowUseCase {
	return &StartWorkflowUseCase{
		workflowRepo:       workflowRepo,
		workflowRunRepo:    workflowRunRepo,
		runWorkflowUseCase: runWorkflowUseCase,
		mercurePublisher:   mercurePublisher,
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

	if len(workflow.Steps) == 0 {
		return errors.New("workflow has no steps")
	}

	err = u.runWorkflowUseCase.Execute(ctx, workflow)
	if err != nil {
		return err
	}

	err = u.mercurePublisher.Publish(fmt.Sprintf("/workflows/%s", workflow.ID),
		map[string]any{
			"type":            "workflow_run.refresh",
			"workflow_run_id": workflow.ID,
			"workflow_id":     workflow.ID,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
