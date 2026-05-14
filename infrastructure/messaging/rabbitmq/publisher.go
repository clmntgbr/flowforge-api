package rabbitmq

import (
	"context"
	"encoding/json"
	"flowforge-api/infrastructure/config"
	"fmt"

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

func NewPublisherFromEnv(env *config.Config) (Publisher, error) {
	conn, err := amqp.Dial(env.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("channel: %w", err)
	}
	return NewPublisher(ch), nil
}

func (p *publisher) PublishStepRunEvent(ctx context.Context, config *config.Config, event StepRunEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	fmt.Println("🔄 Publishing step run event", event)
	fmt.Println("🔄 Exchange name", config.ExchangeName)
	fmt.Println("🔄 Runner queue name", config.RunnerQueueName)
	fmt.Println("🔄 Body", string(body))

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
