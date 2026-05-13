package consumer

import (
	"context"
)

type FailWorkflowStepUseCase struct {
}

func NewFailWorkflowStepUseCase() *FailWorkflowStepUseCase {
	return &FailWorkflowStepUseCase{}
}

func (u *FailWorkflowStepUseCase) Execute(ctx context.Context) error {
	return nil
}
