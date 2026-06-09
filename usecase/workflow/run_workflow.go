package workflow

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow_run"
	"fmt"
)

type RunWorkflowUseCase struct {
	workflowRunRepo           *repository.WorkflowRunRepository
	stepRepo                  *repository.StepRepository
	createWorkflowRunUseCase  *workflow_run.CreateWorkflowRunUseCase
	hasStepRunUseCase         *step_run.HasStepRunUseCase
	createStepRunUseCase      *step_run.CreateStepRunUseCase
	executeStepRunUseCase     *step_run.ExecuteStepRunUseCase
	executeWorkflowRunUseCase *workflow_run.ExecuteWorkflowRunUseCase
	env                       *config.Config
	stepRunPublisher          rabbitmq.Publisher
}

func NewRunWorkflowUseCase(
	workflowRepo *repository.WorkflowRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
	stepRepo *repository.StepRepository,
	createWorkflowRunUseCase *workflow_run.CreateWorkflowRunUseCase,
	hasStepRunUseCase *step_run.HasStepRunUseCase,
	createStepRunUseCase *step_run.CreateStepRunUseCase,
	executeStepRunUseCase *step_run.ExecuteStepRunUseCase,
	executeWorkflowRunUseCase *workflow_run.ExecuteWorkflowRunUseCase,
	env *config.Config,
	stepRunPublisher rabbitmq.Publisher,
) *RunWorkflowUseCase {
	return &RunWorkflowUseCase{
		workflowRunRepo:           workflowRunRepo,
		stepRepo:                  stepRepo,
		createWorkflowRunUseCase:  createWorkflowRunUseCase,
		hasStepRunUseCase:         hasStepRunUseCase,
		createStepRunUseCase:      createStepRunUseCase,
		executeStepRunUseCase:     executeStepRunUseCase,
		executeWorkflowRunUseCase: executeWorkflowRunUseCase,
		env:                       env,
		stepRunPublisher:          stepRunPublisher,
	}
}

func (u *RunWorkflowUseCase) Execute(txCtx context.Context, workflow entity.Workflow) error {
	step, err := (*u.stepRepo).GetFirstStepByWorkflowID(txCtx, workflow.ID)
	if err != nil {
		return fmt.Errorf("🚨 failed to get steps by workflow ID: %w", err)
	}

	if step == nil {
		return nil
	}

	workflowRun, err := (*u.workflowRunRepo).GetByWorkflowIDAndNotEnded(txCtx, workflow.ID)
	if err != nil {
		return fmt.Errorf("🚨 failed to get workflow run by workflow ID and not ended: %w", err)
	}

	if workflowRun == nil {
		workflowRun, err = u.createWorkflowRunUseCase.Execute(txCtx, workflow.ID)
		if err != nil {
			return fmt.Errorf("🚨 failed to create workflow run: %w", err)
		}
	}

	if workflowRun == nil {
		return errors.New("workflow run not found")
	}

	if workflowRun.Status == enum.WorkflowRunStatusRunning {
		return errors.New("workflow run is running")
	}

	hasStepRun := u.hasStepRunUseCase.Execute(txCtx, workflowRun.ID)
	if hasStepRun {
		return errors.New("step run already exists")
	}

	stepRun, err := u.createStepRunUseCase.Execute(txCtx, workflowRun.ID, step.ID)
	if err != nil {
		return fmt.Errorf("🚨 failed to create step run: %w", err)
	}

	stepRun, err = u.executeStepRunUseCase.Execute(txCtx, &stepRun)
	if err != nil {
		return fmt.Errorf("🚨 failed to execute step run: %w", err)
	}

	_, err = u.executeWorkflowRunUseCase.Execute(txCtx, *workflowRun)
	if err != nil {
		return fmt.Errorf("🚨 failed to execute workflow run: %w", err)
	}

	stepRun.Step = *step
	stepRun.WorkflowRun = *workflowRun

	if u.stepRunPublisher == nil {
		return fmt.Errorf("step run publisher is not configured")
	}

	event := rabbitmq.NewStepRunEvent(stepRun)
	if err := u.stepRunPublisher.PublishStepRunEvent(txCtx, u.env, event); err != nil {
		return fmt.Errorf("🚨 failed to publish step run: %w", err)
	}

	return nil
}
