package handler

import (
	"context"
	"flowforge-api/infrastructure/config"
	rabbitmqDTO "flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/infrastructure/messaging/security"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RunnerHandler struct {
	env               *config.Config
	securityValidator *security.WorkerSecurityValidator
	parser            *security.WorkerParser
}

func NewRunnerHandler(env *config.Config) *RunnerHandler {
	return &RunnerHandler{
		env:               env,
		parser:            security.NewWorkerParser(env),
		securityValidator: security.NewWorkerSecurityValidator(env),
	}
}

func (h *RunnerHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	var payload rabbitmqDTO.MessagePayload
	if err := h.parser.ParseAndValidate(message.Body, &payload); err != nil {
		return err
	}

	if err := h.securityValidator.Validate(payload.SecretKey); err != nil {
		return err
	}

	fmt.Println("🔄 Received message", payload)
	return nil
}
