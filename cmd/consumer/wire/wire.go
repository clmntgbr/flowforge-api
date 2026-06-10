package wire

import (
	"flowforge-api/handler"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/mercure"
	rmq "flowforge-api/infrastructure/messaging/rabbitmq"
	repoGorm "flowforge-api/repository/gorm"
	"flowforge-api/usecase/consumer"
	"flowforge-api/usecase/insight"
	usecaseStep "flowforge-api/usecase/step"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/workflow_run"
	"log"

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

	mercurePublisher := mercure.NewPublisher(env.MercureURL, env.MercurePublisherJWTKey)

	createInsightUseCase := insight.NewCreateInsightUseCase(&insightRepo)

	createStepRunUseCase := step_run.NewCreateStepRunUseCase(&stepRunRepo, &stepRepo)
	executeStepRunUseCase := step_run.NewExecuteStepRunUseCase(&stepRunRepo, &stepRepo)
	stepRunPublisher, err := rmq.NewPublisherFromEnv(env)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ publisher: %v", err)
	}
	findNextStepUseCase := usecaseStep.NewFindNextStepUseCase(&stepRepo)
	isCanceledWorkflowRunUseCase := workflow_run.NewIsCanceledWorkflowRunUseCase(&workflowRunRepo)

	failedStepUseCase := consumer.NewFailedStepUseCase(
		createInsightUseCase,
		&stepRunRepo,
		&workflowRunRepo,
		&stepRepo,
		createStepRunUseCase,
		executeStepRunUseCase,
		isCanceledWorkflowRunUseCase,
		stepRunPublisher,
		env,
	)

	completedStepUseCase := consumer.NewCompletedStepUseCase(
		createInsightUseCase,
		&stepRunRepo,
		&workflowRunRepo,
		findNextStepUseCase,
		createStepRunUseCase,
		executeStepRunUseCase,
		isCanceledWorkflowRunUseCase,
		stepRunPublisher,
		env,
	)

	consumerHandler := handler.NewConsumerHandler(
		env,
		completedStepUseCase,
		failedStepUseCase,
		mercurePublisher,
	)

	return &Container{
		ConsumerHandler: consumerHandler,
	}
}
