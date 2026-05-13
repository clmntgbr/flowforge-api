package consumer

import (
	"context"
	consumerDTO "flowforge-api/infrastructure/consumer"
)

type FailWorkflowStepUseCase struct {
}

func NewFailWorkflowStepUseCase() *FailWorkflowStepUseCase {
	return &FailWorkflowStepUseCase{}
}

func (u *FailWorkflowStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerFailedMessage) error {
	return nil
}
