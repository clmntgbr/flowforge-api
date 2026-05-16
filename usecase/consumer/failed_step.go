package consumer

import (
	"context"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"fmt"
)

type FailedStepUseCase struct {
}

func NewFailedStepUseCase() *FailedStepUseCase {
	return &FailedStepUseCase{}
}

func (u *FailedStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerFailedMessage) error {
	fmt.Println(message)
	return nil
}
