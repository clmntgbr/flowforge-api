package handler

import (
	"flowforge-api/usecase/consumer"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler struct {
	completeWorkflowStepUseCase *consumer.CompleteWorkflowStepUseCase
	failWorkflowStepUseCase     *consumer.FailWorkflowStepUseCase
}

func NewConsumerHandler(
	completeWorkflowStepUseCase *consumer.CompleteWorkflowStepUseCase,
	failWorkflowStepUseCase *consumer.FailWorkflowStepUseCase,
) *ConsumerHandler {
	return &ConsumerHandler{
		completeWorkflowStepUseCase: completeWorkflowStepUseCase,
		failWorkflowStepUseCase:     failWorkflowStepUseCase,
	}
}

func (h *ConsumerHandler) HandleMessage(message *amqp.Delivery) error {
	log.Println("received message: ", message.RoutingKey)
	log.Println("message body: ", string(message.Body))
	return nil
}
