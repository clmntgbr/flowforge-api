package workflow

import (
	"context"
	"errors"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"time"

	"github.com/google/uuid"
)

type StopWorkflowUseCase struct {
	workflowRepo       *repository.WorkflowRepository
	workflowRunRepo    *repository.WorkflowRunRepository
	stepRunRepo        *repository.StepRunRepository
	runWorkflowUseCase *RunWorkflowUseCase
}

func NewStopWorkflowUseCase(
	workflowRepo *repository.WorkflowRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
	stepRunRepo *repository.StepRunRepository,
	runWorkflowUseCase *RunWorkflowUseCase,
) *StopWorkflowUseCase {
	return &StopWorkflowUseCase{
		workflowRepo:       workflowRepo,
		workflowRunRepo:    workflowRunRepo,
		stepRunRepo:        stepRunRepo,
		runWorkflowUseCase: runWorkflowUseCase,
	}
}

func (u *StopWorkflowUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) error {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return err
	}

	if workflow.Status == enum.WorkflowStatusInactive {
		return errors.New("workflow is inactive")
	}

	workflowRun, err := (*u.workflowRunRepo).GetByWorkflowIDAndNotEnded(ctx, workflowID)
	if err != nil {
		return err
	}

	if workflowRun.Status != enum.WorkflowRunStatusRunning {
		return errors.New("workflow run is not running")
	}

	workflowRun.Status = enum.WorkflowRunStatusCanceled
	workflowRun.Statuses = append(workflowRun.Statuses, enum.WorkflowRunStatusCanceled)
	now := time.Now().UTC()
	workflowRun.CanceledAt = &now
	workflowRun.CompletedAt = &now

	err = (*u.workflowRunRepo).Update(ctx, workflowRun)
	if err != nil {
		return err
	}

	err = (*u.stepRunRepo).CancelRunningByWorkflowRunID(ctx, workflowRun.ID)
	if err != nil {
		return err
	}

	return nil
}
