package wire

import (
	"flowforge-api/infrastructure/config"
	rmq "flowforge-api/infrastructure/messaging/rabbitmq"
	repoGorm "flowforge-api/repository/gorm"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow"
	"flowforge-api/usecase/workflow_run"

	"gorm.io/gorm"
)

type Container struct {
	ExecuteWorkflowUseCase *workflow.ExecuteWorkflowUseCase
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	stepRepo := repoGorm.NewStepRepository(db)
	workflowRepo := repoGorm.NewWorkflowRepository(db)
	workflowRunRepo := repoGorm.NewWorkflowRunRepository(db)
	stepRunRepo := repoGorm.NewStepRunRepository(db)

	createWorkflowRunUseCase := workflow_run.NewCreateWorkflowRunUseCase(&workflowRunRepo)
	hasStepRunUseCase := step_run.NewHasStepRunUseCase(&stepRunRepo)
	createStepRunUseCase := step_run.NewCreateStepRunUseCase(&stepRunRepo, &stepRepo)
	executeStepRunUseCase := step_run.NewExecuteStepRunUseCase(&stepRunRepo, &stepRepo)
	executeWorkflowRunUseCase := workflow_run.NewExecuteWorkflowRunUseCase(&workflowRunRepo)
	stepRunPublisher := rmq.NewPublisherFromEnv(env)

	executeWorkflowUseCase := workflow.NewExecuteWorkflowUseCase(
		&workflowRepo,
		&workflowRunRepo,
		&stepRepo,
		createWorkflowRunUseCase,
		hasStepRunUseCase,
		createStepRunUseCase,
		executeStepRunUseCase,
		executeWorkflowRunUseCase,
		env,
		&stepRunPublisher,
	)

	return &Container{
		ExecuteWorkflowUseCase: executeWorkflowUseCase,
	}
}
