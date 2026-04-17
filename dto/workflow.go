package dto

import (
	"forgeflow-api/domain"
	"time"
)

type WorkflowOutput struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateWorkflowInput struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty,min=2,max=255"`
}

func NewWorkflowOutput(workflow domain.Workflow) WorkflowOutput {
	return WorkflowOutput{
		ID:          workflow.ID.String(),
		Name:        workflow.Name,
		Description: workflow.Description,
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
	}
}

func NewWorkflowsOutput(workflows []domain.Workflow) []WorkflowOutput {
	outputs := make([]WorkflowOutput, len(workflows))
	for i, workflow := range workflows {
		outputs[i] = NewWorkflowOutput(workflow)
	}
	return outputs
}
