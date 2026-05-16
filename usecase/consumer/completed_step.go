package consumer

import (
	"context"
	consumerDTO "flowforge-api/infrastructure/consumer"
	"fmt"
)

type CompletedStepUseCase struct {
}

func NewCompletedStepUseCase() *CompletedStepUseCase {
	return &CompletedStepUseCase{}
}

func (u *CompletedStepUseCase) Execute(ctx context.Context, message consumerDTO.ConsumerCompletedMessage) error {
	fmt.Println(message)
	return nil
}
