package wire

import (
	"flowforge-api/infrastructure/config"
	rmq "flowforge-api/infrastructure/messaging/rabbitmq"
	repoGorm "flowforge-api/repository/gorm"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow"
	"flowforge-api/usecase/workflow_run"
	"log"

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

	stepRunPublisher, err := rmq.NewPublisherFromEnv(env)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ publisher: %v", err)
	}

	createWorkflowRunUseCase := workflow_run.NewCreateWorkflowRunUseCase(&workflowRunRepo)
	hasStepRunUseCase := step_run.NewHasStepRunUseCase(&stepRunRepo)
	createStepRunUseCase := step_run.NewCreateStepRunUseCase(&stepRunRepo, &stepRepo)
	executeStepRunUseCase := step_run.NewExecuteStepRunUseCase(&stepRunRepo, &stepRepo)
	executeWorkflowRunUseCase := workflow_run.NewExecuteWorkflowRunUseCase(&workflowRunRepo)
	runWorkflowUseCase := workflow.NewRunWorkflowUseCase(&workflowRepo, &workflowRunRepo, &stepRepo, createWorkflowRunUseCase, hasStepRunUseCase, createStepRunUseCase, executeStepRunUseCase, executeWorkflowRunUseCase, env, stepRunPublisher)

	executeWorkflowUseCase := workflow.NewExecuteWorkflowUseCase(
		&workflowRepo,
		runWorkflowUseCase,
		env,
	)

	return &Container{
		ExecuteWorkflowUseCase: executeWorkflowUseCase,
	}
}
