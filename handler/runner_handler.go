package handler

import (
	"context"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	rabbitmqDTO "flowforge-api/infrastructure/messaging/rabbitmq"
	"flowforge-api/infrastructure/messaging/security"
	"flowforge-api/infrastructure/runner"
	"flowforge-api/usecase/step"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RunnerHandler struct {
	env               *config.Config
	securityValidator *security.WorkerSecurityValidator
	parser            *security.WorkerParser
	runStepUseCase    *step.RunStepUseCase
	publisher         rabbitmq.Publisher
}

func NewRunnerHandler(env *config.Config, runStepUseCase *step.RunStepUseCase, publisher rabbitmq.Publisher) *RunnerHandler {
	return &RunnerHandler{
		env:               env,
		parser:            security.NewWorkerParser(env),
		securityValidator: security.NewWorkerSecurityValidator(env),
		runStepUseCase:    runStepUseCase,
		publisher:         publisher,
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

	response, err := h.runStepUseCase.Execute(ctx, &payload.StepRunEvent)
	if err != nil {
		return h.PublishFailure(ctx, payload.StepRunEvent, response, err)
	}

	fmt.Println("🔄 Received message", payload)
	return h.PublishSuccess(ctx, payload.StepRunEvent, response)
}

func (h *RunnerHandler) PublishSuccess(ctx context.Context, event rabbitmqDTO.StepRunEvent, response runner.RunnerResponse) error {
	message := rabbitmqDTO.RunnerCompletedMessage{
		WorkflowRunID: event.WorkflowRunID.String(),
		StepRunID:     event.StepRunID.String(),
		CompletedAt:   time.Now().Format(time.RFC3339),
		Insights:      response.Insights,
		Response:      response.Response,
	}

	if err := h.publisher.PublishRunnerCompleted(ctx, h.env, message); err != nil {
		log.Printf("client_handler: failed to publish success message (step_run_id=%s): %v", event.StepRunID, err)
		return fmt.Errorf("failed to publish success message: %w", err)
	}

	return nil
}

func (h *RunnerHandler) PublishFailure(ctx context.Context, event rabbitmqDTO.StepRunEvent, response runner.RunnerResponse, execError error) error {
	message := rabbitmqDTO.RunnerFailedMessage{
		WorkflowRunID: event.WorkflowRunID.String(),
		StepRunID:     event.StepRunID.String(),
		Error:         execError.Error(),
		FailedAt:      time.Now().Format(time.RFC3339),
		Insights:      response.Insights,
		Response:      response.Response,
	}

	if err := h.publisher.PublishRunnerFailed(ctx, h.env, message); err != nil {
		log.Printf("client_handler: failed to publish failure message (step_run_id=%s): %v", event.StepRunID, err)
		return fmt.Errorf("failed to publish failure message (original error: %v): %w", execError, err)
	}

	return fmt.Errorf("step execution failed: %w", execError)
}
