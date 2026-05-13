package consumer

import (
	"context"
	consumerDTO "flowforge-api/infrastructure/consumer"
)

type CompleteWorkflowStepUseCase struct {
}

func NewCompleteWorkflowStepUseCase() *CompleteWorkflowStepUseCase {
	return &CompleteWorkflowStepUseCase{}
}

func (u *CompleteWorkflowStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerCompletedMessage) error {
	return nil
}
