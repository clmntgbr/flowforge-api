package handler

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RunnerHandler struct {
}

func NewRunnerHandler() *RunnerHandler {
	return &RunnerHandler{}
}

func (h *RunnerHandler) HandleMessage(message *amqp.Delivery) error {
	log.Println("received message: ", message.RoutingKey)
	log.Println("message body: ", string(message.Body))
	return nil
}
