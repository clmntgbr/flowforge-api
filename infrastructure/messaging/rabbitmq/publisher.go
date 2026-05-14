package rabbitmq

import (
	"context"
	"encoding/json"
	"flowforge-api/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	PublishStepRunEvent(ctx context.Context, config *config.Config, event StepRunEvent) error
}

type publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) Publisher {
	return &publisher{
		channel: channel,
	}
}

func (p *publisher) PublishStepRunEvent(ctx context.Context, config *config.Config, event StepRunEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		config.ExchangeName,
		config.RunnerQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
