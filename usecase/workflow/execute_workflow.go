package workflow

import (
	"context"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow_run"
	"fmt"
	"log"
)

type ExecuteWorkflowUseCase struct {
	workflowRepo              repository.WorkflowRepository
	workflowRunRepo           repository.WorkflowRunRepository
	stepRepo                  repository.StepRepository
	createWorkflowRunUseCase  *workflow_run.CreateWorkflowRunUseCase
	hasStepRunUseCase         *step_run.HasStepRunUseCase
	createStepRunUseCase      *step_run.CreateStepRunUseCase
	executeStepRunUseCase     *step_run.ExecuteStepRunUseCase
	executeWorkflowRunUseCase *workflow_run.ExecuteWorkflowRunUseCase
}

func NewExecuteWorkflowUseCase(
	workflowRepo repository.WorkflowRepository,
	workflowRunRepo repository.WorkflowRunRepository,
	stepRepo repository.StepRepository,
	createWorkflowRunUseCase *workflow_run.CreateWorkflowRunUseCase,
	hasStepRunUseCase *step_run.HasStepRunUseCase,
	createStepRunUseCase *step_run.CreateStepRunUseCase,
	executeStepRunUseCase *step_run.ExecuteStepRunUseCase,
	executeWorkflowRunUseCase *workflow_run.ExecuteWorkflowRunUseCase,
) *ExecuteWorkflowUseCase {
	return &ExecuteWorkflowUseCase{
		workflowRepo:              workflowRepo,
		workflowRunRepo:           workflowRunRepo,
		stepRepo:                  stepRepo,
		createWorkflowRunUseCase:  createWorkflowRunUseCase,
		hasStepRunUseCase:         hasStepRunUseCase,
		createStepRunUseCase:      createStepRunUseCase,
		executeStepRunUseCase:     executeStepRunUseCase,
		executeWorkflowRunUseCase: executeWorkflowRunUseCase,
	}
}

func (u *ExecuteWorkflowUseCase) Execute(ctx context.Context) error {
	log.Println("🔄 Executing workflow use case")

	workflows, err := u.workflowRepo.GetWorkflowsForExecution(ctx)
	if err != nil {
		return fmt.Errorf("🚨 failed to get workflows for execution: %w", err)
	}

	for _, workflow := range workflows {
		log.Println("🔄 Executing workflow", workflow.ID)

		workflowRun, err := u.workflowRunRepo.GetByWorkflowIDAndNotEnded(ctx, workflow.ID)
		if err != nil {
			return fmt.Errorf("🚨 failed to get workflow run by workflow ID and not ended: %w", err)
		}

		if workflowRun == nil {
			if workflowRun, err = u.createWorkflowRunUseCase.Execute(ctx, workflow.ID); err != nil {
				return fmt.Errorf("🚨 failed to create workflow run: %w", err)
			}
		}

		if workflowRun == nil {
			continue
		}

		fmt.Println("🔄 Workflow run", workflowRun)
		fmt.Println("🔄 Workflow run status", workflowRun.Status)
		if workflowRun.Status == enum.WorkflowRunStatusRunning {
			continue
		}

		step, err := u.stepRepo.GetFirstStepByWorkflowID(ctx, workflow.ID)
		if err != nil {
			return fmt.Errorf("🚨 failed to get steps by workflow ID: %w", err)
		}

		if step == nil {
			continue
		}

		hasStepRun, err := u.hasStepRunUseCase.Execute(ctx, workflowRun.ID)
		if err != nil {
			return fmt.Errorf("🚨 failed to check if step run exists: %w", err)
		}

		if hasStepRun {
			continue
		}

		stepRun, err := u.createStepRunUseCase.Execute(ctx, workflowRun.ID, step.ID)
		if err != nil {
			return fmt.Errorf("🚨 failed to create step run: %w", err)
		}

		stepRun, err = u.executeStepRunUseCase.Execute(ctx, &stepRun)
		if err != nil {
			return fmt.Errorf("🚨 failed to execute step run: %w", err)
		}

		workflowRun, err = u.executeWorkflowRunUseCase.Execute(ctx, workflowRun)
		if err != nil {
			return fmt.Errorf("🚨 failed to execute workflow run: %w", err)
		}

		fmt.Println("🔄 Step", step)
		log.Println("🔄 Creating workflow run", workflow.ID)
	}

	return nil
}
