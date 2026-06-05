package wire

import (
	"flowforge-api/handler"
	"flowforge-api/infrastructure/config"
	rmq "flowforge-api/infrastructure/messaging/rabbitmq"
	repoGorm "flowforge-api/repository/gorm"
	"flowforge-api/usecase/consumer"
	"flowforge-api/usecase/insight"
	"flowforge-api/usecase/step_run"

	"gorm.io/gorm"
)

type Container struct {
	ConsumerHandler *handler.ConsumerHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	stepRepo := repoGorm.NewStepRepository(db)
	workflowRunRepo := repoGorm.NewWorkflowRunRepository(db)
	stepRunRepo := repoGorm.NewStepRunRepository(db)
	insightRepo := repoGorm.NewInsightRepository(db)

	createInsightUseCase := insight.NewCreateInsightUseCase(&insightRepo)

	createStepRunUseCase := step_run.NewCreateStepRunUseCase(&stepRunRepo, &stepRepo)
	executeStepRunUseCase := step_run.NewExecuteStepRunUseCase(&stepRunRepo, &stepRepo)
	stepRunPublisher := rmq.NewPublisherFromEnv(env)

	failedStepUseCase := consumer.NewFailedStepUseCase(
		createInsightUseCase,
		&stepRunRepo,
		&workflowRunRepo,
		&stepRepo,
		createStepRunUseCase,
		executeStepRunUseCase,
		&stepRunPublisher,
		env,
	)

	completedStepUseCase := consumer.NewCompletedStepUseCase(
		createInsightUseCase,
		&stepRunRepo,
		&workflowRunRepo,
		&stepRepo,
		createStepRunUseCase,
		executeStepRunUseCase,
		&stepRunPublisher,
		env,
	)

	consumerHandler := handler.NewConsumerHandler(
		env,
		completedStepUseCase,
		failedStepUseCase,
	)

	return &Container{
		ConsumerHandler: consumerHandler,
	}
}
