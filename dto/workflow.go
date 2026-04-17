package dto

import (
	"forgeflow-api/domain"
	"time"
)

type MinimalWorkflowOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WorkflowOutput struct {
	MinimalWorkflowOutput
	Description string `json:"description"`
}

type CreateWorkflowInput struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty,min=2,max=255"`
}

type UpdateWorkflowInput struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty,min=2,max=255"`
}

func NewMinimalWorkflowOutput(workflow domain.Workflow) MinimalWorkflowOutput {
	return MinimalWorkflowOutput{
		ID:        workflow.ID.String(),
		Name:      workflow.Name,
		CreatedAt: workflow.CreatedAt,
		UpdatedAt: workflow.UpdatedAt,
	}
}

func NewWorkflowOutput(workflow domain.Workflow) WorkflowOutput {
	return WorkflowOutput{
		MinimalWorkflowOutput: NewMinimalWorkflowOutput(workflow),
		Description:           workflow.Description,
	}
}

func NewMinimalWorkflowsOutput(workflows []domain.Workflow) []MinimalWorkflowOutput {
	outputs := make([]MinimalWorkflowOutput, len(workflows))
	for i, workflow := range workflows {
		outputs[i] = NewMinimalWorkflowOutput(workflow)
	}
	return outputs
}
