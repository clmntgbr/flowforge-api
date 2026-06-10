package wire

import (
	"flowforge-api/handler"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/mercure"
	rmq "flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/usecase/step"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type Container struct {
	RunnerHandler *handler.RunnerHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	_ = db

	runStepUseCase := step.NewRunStepUseCase(http.DefaultClient)

	mercurePublisher := mercure.NewPublisher(env.MercureURL, env.MercurePublisherJWTKey)

	stepRunPublisher, err := rmq.NewPublisherFromEnv(env)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ publisher: %v", err)
	}

	runnerHandler := handler.NewRunnerHandler(env, runStepUseCase, stepRunPublisher, mercurePublisher)

	return &Container{
		RunnerHandler: runnerHandler,
	}
}
