package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"flowforge-api/infrastructure/config"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/infrastructure/messaging/security"
	"flowforge-api/usecase/consumer"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler struct {
	env                  *config.Config
	completedStepUseCase *consumer.CompletedStepUseCase
	failedStepUseCase    *consumer.FailedStepUseCase
	securityValidator    *security.WorkerSecurityValidator
	parser               *security.WorkerParser
}

func NewConsumerHandler(
	env *config.Config,
	completedStepUseCase *consumer.CompletedStepUseCase,
	failedStepUseCase *consumer.FailedStepUseCase,
) *ConsumerHandler {
	return &ConsumerHandler{
		env:                  env,
		completedStepUseCase: completedStepUseCase,
		failedStepUseCase:    failedStepUseCase,
		parser:               security.NewWorkerParser(env),
		securityValidator:    security.NewWorkerSecurityValidator(env),
	}
}

func (h *ConsumerHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	fmt.Println("🔄 Handling message", message.RoutingKey)
	switch message.RoutingKey {
	case h.env.ConsumerRoutingKeyCompleted:
		var payload consumerDTO.ConsumerMessage
		if err := h.parser.ParseAndValidate(message.Body, &payload); err != nil {
			fmt.Println("🚨 Error parsing message", err)
			return err
		}
		if err := h.securityValidator.Validate(payload.SecretKey); err != nil {
			fmt.Println("🚨 Error validating message", err)
			return err
		}

		var completed consumerDTO.ConsumerCompletedMessage
		if err := json.Unmarshal(payload.Message, &completed); err != nil {
			return fmt.Errorf("decode completed message: %w", err)
		}

		if err := h.completedStepUseCase.Execute(ctx, completed); err != nil {
			fmt.Println("🚨 Error executing complete workflow step", err)
			return err
		}

	case h.env.ConsumerRoutingKeyFailed:
		var payload consumerDTO.ConsumerMessage
		if err := h.parser.ParseAndValidate(message.Body, &payload); err != nil {
			fmt.Println("🚨 Error parsing message", err)
			return err
		}
		if err := h.securityValidator.Validate(payload.SecretKey); err != nil {
			fmt.Println("🚨 Error validating message", err)
			return err
		}
		var failed consumerDTO.ConsumerFailedMessage
		if err := json.Unmarshal(payload.Message, &failed); err != nil {
			return fmt.Errorf("decode failed message: %w", err)
		}

		err := h.failedStepUseCase.Execute(ctx, failed)
		if err != nil {
			fmt.Println("🚨 Error executing fail workflow step", err)
			return err
		}

	default:
		return fmt.Errorf("unsupported routing key %q", message.RoutingKey)
	}

	fmt.Println("✅ Success handling message", message.RoutingKey)
	return nil
}
