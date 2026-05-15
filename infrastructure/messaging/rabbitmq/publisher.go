package rabbitmq

import (
	"context"
	"encoding/json"
	"flowforge-api/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	PublishStepRunEvent(ctx context.Context, config *config.Config, event StepRunEvent) error
	PublishRunnerCompleted(ctx context.Context, config *config.Config, message RunnerCompletedMessage) error
	PublishRunnerFailed(ctx context.Context, config *config.Config, message RunnerFailedMessage) error
}

type publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) Publisher {
	return &publisher{
		channel: channel,
	}
}

func NewPublisherFromEnv(env *config.Config) Publisher {
	conn, err := amqp.Dial(env.RabbitMQURL)
	if err != nil {
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil
	}
	return NewPublisher(ch)
}

func (p *publisher) PublishStepRunEvent(ctx context.Context, config *config.Config, event StepRunEvent) error {
	message := MessagePayload{
		SecretKey:    config.RabbitMQSecretKey,
		StepRunEvent: event,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		"",
		config.RunnerQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *publisher) PublishRunnerCompleted(ctx context.Context, config *config.Config, message RunnerCompletedMessage) error {
	response := MessageResponse{
		SecretKey: config.RabbitMQSecretKey,
		Message:   message,
	}

	body, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		config.ExchangeName,
		config.ConsumerRoutingKeyCompleted,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *publisher) PublishRunnerFailed(ctx context.Context, config *config.Config, message RunnerFailedMessage) error {
	response := MessageResponse{
		SecretKey: config.RabbitMQSecretKey,
		Message:   message,
	}

	body, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		config.ExchangeName,
		config.ConsumerRoutingKeyFailed,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
