package consumer

import (
	"context"
)

type CompleteWorkflowStepUseCase struct {
}

func NewCompleteWorkflowStepUseCase() *CompleteWorkflowStepUseCase {
	return &CompleteWorkflowStepUseCase{}
}

func (u *CompleteWorkflowStepUseCase) Execute(ctx context.Context) error {
	return nil
}
