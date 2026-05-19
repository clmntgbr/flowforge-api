package wire

import (
	"flowforge-api/handler"
	"flowforge-api/infrastructure/config"
	rmq "flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/usecase/step"
	"net/http"

	"gorm.io/gorm"
)

type Container struct {
	RunnerHandler *handler.RunnerHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	_ = db

	runStepUseCase := step.NewRunStepUseCase(http.DefaultClient)
	stepRunPublisher := rmq.NewPublisherFromEnv(env)
	runnerHandler := handler.NewRunnerHandler(env, runStepUseCase, &stepRunPublisher)

	return &Container{
		RunnerHandler: runnerHandler,
	}
}
