package handler

import amqp "github.com/rabbitmq/amqp091-go"

type ConsumerHandler struct {
}

func NewConsumerHandler() *ConsumerHandler {
	return &ConsumerHandler{}
}

func (h *ConsumerHandler) HandleMessage(message *amqp.Delivery) error {
	return nil
}
