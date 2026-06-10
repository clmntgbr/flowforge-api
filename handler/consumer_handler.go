package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"flowforge-api/infrastructure/config"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"flowforge-api/infrastructure/mercure"
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
	mercurePublisher     *mercure.Publisher
}

func NewConsumerHandler(
	env *config.Config,
	completedStepUseCase *consumer.CompletedStepUseCase,
	failedStepUseCase *consumer.FailedStepUseCase,
	mercurePublisher *mercure.Publisher,
) *ConsumerHandler {
	return &ConsumerHandler{
		env:                  env,
		completedStepUseCase: completedStepUseCase,
		failedStepUseCase:    failedStepUseCase,
		parser:               security.NewWorkerParser(env),
		securityValidator:    security.NewWorkerSecurityValidator(env),
		mercurePublisher:     mercurePublisher,
	}
}

func (h *ConsumerHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	fmt.Println("🔄 Handling message", message.RoutingKey)
	var payload consumerDTO.ConsumerMessage
	switch message.RoutingKey {
	case h.env.ConsumerRoutingKeyCompleted:
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

		fmt.Printf("🚨 publishing to /workflows/%s\n", completed.WorkflowID)
		err := h.mercurePublisher.Publish(fmt.Sprintf("/workflows/%s", completed.WorkflowID),
			map[string]any{
				"type":            "workflow_run.refresh",
				"workflow_run_id": completed.WorkflowRunID,
				"workflow_id":     completed.WorkflowID,
			},
		)

		if err != nil {
			return fmt.Errorf("🚨 failed to publish workflow run failed: %w", err)
		}

	case h.env.ConsumerRoutingKeyFailed:
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

		if err := h.failedStepUseCase.Execute(ctx, failed); err != nil {
			fmt.Println("🚨 Error executing failed workflow step", err)
			return err
		}

		fmt.Printf("🚨 publishing to /workflows/%s\n", failed.WorkflowID)
		err := h.mercurePublisher.Publish(fmt.Sprintf("/workflows/%s", failed.WorkflowID),
			map[string]any{
				"type":            "workflow_run.refresh",
				"workflow_run_id": failed.WorkflowRunID,
				"workflow_id":     failed.WorkflowID,
			},
		)

		if err != nil {
			return fmt.Errorf("🚨 failed to publish workflow run failed: %w", err)
		}

	default:
		return fmt.Errorf("unsupported routing key %q", message.RoutingKey)
	}

	fmt.Println("✅ Success handling message", message.RoutingKey)
	return nil
}
