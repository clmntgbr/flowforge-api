package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"flowforge-api/infrastructure/config"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/usecase/consumer"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler struct {
	env                         *config.Config
	completeWorkflowStepUseCase *consumer.CompleteWorkflowStepUseCase
	failWorkflowStepUseCase     *consumer.FailWorkflowStepUseCase
	securityValidator           *rabbitmq.WorkerSecurityValidator
	parser                      *rabbitmq.WorkerParser
}

func NewConsumerHandler(
	env *config.Config,
	completeWorkflowStepUseCase *consumer.CompleteWorkflowStepUseCase,
	failWorkflowStepUseCase *consumer.FailWorkflowStepUseCase,
) *ConsumerHandler {
	return &ConsumerHandler{
		env:                         env,
		completeWorkflowStepUseCase: completeWorkflowStepUseCase,
		failWorkflowStepUseCase:     failWorkflowStepUseCase,
		parser:                      rabbitmq.NewWorkerParser(env),
		securityValidator:           rabbitmq.NewWorkerSecurityValidator(env),
	}
}

func (h *ConsumerHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	switch message.RoutingKey {
	case h.env.ConsumerRoutingKeyCompleted:
		var payload consumerDTO.ConsumerMessage
		if err := h.parser.ParseAndValidate(message.Body, &payload); err != nil {
			return err
		}
		if err := h.securityValidator.Validate(payload.SecretKey); err != nil {
			return err
		}
		var completed consumerDTO.ConsumerCompletedMessage
		if err := json.Unmarshal(payload.Message, &completed); err != nil {
			return fmt.Errorf("decode completed message: %w", err)
		}
		if err := h.completeWorkflowStepUseCase.Execute(ctx, completed); err != nil {
			return err
		}

	case h.env.ConsumerRoutingKeyFailed:
		var payload consumerDTO.ConsumerMessage
		if err := h.parser.ParseAndValidate(message.Body, &payload); err != nil {
			return err
		}
		if err := h.securityValidator.Validate(payload.SecretKey); err != nil {
			return err
		}
		var failed consumerDTO.ConsumerFailedMessage
		if err := json.Unmarshal(payload.Message, &failed); err != nil {
			return fmt.Errorf("decode failed message: %w", err)
		}
		if err := h.failWorkflowStepUseCase.Execute(ctx, failed); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported routing key %q", message.RoutingKey)
	}

	return nil
}
