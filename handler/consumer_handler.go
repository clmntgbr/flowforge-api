package handler

import (
	"context"
	"flowforge-api/infrastructure/config"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/usecase/consumer"
	"fmt"

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
		if err := h.completeWorkflowStepUseCase.Execute(payload.Message); err != nil {
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
		if err := h.failWorkflowStepUseCase.Execute(payload.Message); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported routing key %q", message.RoutingKey)
	}

	return nil
}
