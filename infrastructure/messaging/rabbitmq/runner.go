package rabbitmq

import (
	"fmt"
	"log"

	"flowforge-api/handler"
	"flowforge-api/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type WorkerRunner struct {
	env     *config.Config
	handler *handler.RunnerHandler

	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewWorkerRunner(
	env *config.Config,
	handler *handler.RunnerHandler,
) *WorkerRunner {
	return &WorkerRunner{
		env:     env,
		handler: handler,
	}
}

func (c *WorkerRunner) Start() error {
	conn, err := amqp.Dial(c.env.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	c.channel = channel

	if err := c.channel.ExchangeDeclare(
		c.env.ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf(
			"failed to declare exchange %q: %w",
			c.env.ExchangeName,
			err,
		)
	}

	queue, err := c.channel.QueueDeclare(
		c.env.RunnerQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare queue %q: %w",
			c.env.RunnerQueueName,
			err,
		)
	}

	messages, err := c.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to consume queue %q: %w",
			queue.Name,
			err,
		)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Printf(
		"[*] Waiting for messages on exchange %q",
		c.env.ExchangeName,
	)

	for message := range messages {
		if err := c.handler.HandleMessage(&message); err != nil {
			log.Printf(
				"rejected message (routing key: %q): %v",
				message.RoutingKey,
				err,
			)

			if nackErr := message.Nack(false, false); nackErr != nil {
				log.Printf("failed to nack message: %v", nackErr)
			}

			continue
		}

		if ackErr := message.Ack(false); ackErr != nil {
			log.Printf("failed to ack message: %v", ackErr)
		}
	}

	return nil
}

func (c *WorkerRunner) Stop() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
