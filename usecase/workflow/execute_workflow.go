package workflow

import (
	"context"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/config"
	repogorm "flowforge-api/repository/gorm"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type ExecuteWorkflowUseCase struct {
	workflowRepo       *repository.WorkflowRepository
	runWorkflowUseCase *RunWorkflowUseCase
	env                *config.Config
}

func NewExecuteWorkflowUseCase(
	workflowRepo *repository.WorkflowRepository,
	runWorkflowUseCase *RunWorkflowUseCase,
	env *config.Config,
) *ExecuteWorkflowUseCase {
	return &ExecuteWorkflowUseCase{
		workflowRepo:       workflowRepo,
		runWorkflowUseCase: runWorkflowUseCase,
		env:                env,
	}
}

func (u *ExecuteWorkflowUseCase) Execute(ctx context.Context) error {
	log.Println("🔄 Executing workflow use case")

	workflows, err := (*u.workflowRepo).GetWorkflowsForExecution(ctx)
	if err != nil {
		return fmt.Errorf("🚨 failed to get workflows for execution: %w", err)
	}

	for _, wf := range workflows {
		err := (*u.workflowRepo).Transaction(ctx, func(tx *gorm.DB) error {
			txCtx := repogorm.ContextWithTx(ctx, tx)
			return u.runWorkflowUseCase.Execute(txCtx, wf)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
