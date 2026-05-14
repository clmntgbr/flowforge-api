package workflow

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	repogorm "flowforge-api/repository/gorm"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow_run"
	"fmt"
	"log"

	"gorm.io/gorm"
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
	env                       *config.Config
	stepRunPublisher          rabbitmq.Publisher
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
	env *config.Config,
	stepRunPublisher rabbitmq.Publisher,
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
		env:                       env,
		stepRunPublisher:          stepRunPublisher,
	}
}

func (u *ExecuteWorkflowUseCase) Execute(ctx context.Context) error {
	log.Println("🔄 Executing workflow use case")

	workflows, err := u.workflowRepo.GetWorkflowsForExecution(ctx)
	if err != nil {
		return fmt.Errorf("🚨 failed to get workflows for execution: %w", err)
	}

	for _, wf := range workflows {
		wf := wf
		err := u.workflowRepo.Transaction(ctx, func(tx *gorm.DB) error {
			txCtx := repogorm.ContextWithTx(ctx, tx)
			return u.runExecuteWorkflowIteration(txCtx, wf)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *ExecuteWorkflowUseCase) runExecuteWorkflowIteration(txCtx context.Context, workflow entity.Workflow) error {
	log.Println("🔄 Executing workflow", workflow.ID)

	workflowRun, err := u.workflowRunRepo.GetByWorkflowIDAndNotEnded(txCtx, workflow.ID)
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
		return nil
	}

	fmt.Println("🔄 Workflow run", workflowRun)
	fmt.Println("🔄 Workflow run status", workflowRun.Status)
	if workflowRun.Status == enum.WorkflowRunStatusRunning {
		return nil
	}

	step, err := u.stepRepo.GetFirstStepByWorkflowID(txCtx, workflow.ID)
	if err != nil {
		return fmt.Errorf("🚨 failed to get steps by workflow ID: %w", err)
	}

	if step == nil {
		return nil
	}

	hasStepRun := u.hasStepRunUseCase.Execute(txCtx, workflowRun.ID)
	if hasStepRun {
		fmt.Println("🔄 Step run already exists")
		return nil
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
		return fmt.Errorf("🚨 step run publisher is not configured")
	}

	event := rabbitmq.NewStepRunEvent(stepRun)
	if err := u.stepRunPublisher.PublishStepRunEvent(txCtx, u.env, event); err != nil {
		return fmt.Errorf("🚨 failed to publish step run: %w", err)
	}

	fmt.Println("🔄 Step", step)
	log.Println("🔄 Creating workflow run", workflow.ID)

	return nil
}
