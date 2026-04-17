package dto

import (
	"forgeflow-api/domain"
	"time"

	"github.com/google/uuid"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type StepOutput struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Position    Position        `json:"position"`
	Index       string          `json:"index"`
	Timeout     int             `json:"timeout"`
	Endpoint    *EndpointOutput `json:"endpoint,omitempty"`
	EndpointID  string          `json:"endpointId"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

type UpdateWorkflowStepInput struct {
	ID         string   `json:"id" validate:"required,uuid"`
	Name       string   `json:"name" validate:"required,min=2,max=100"`
	EndpointID string   `json:"endpointId" validate:"required,uuid"`
	Position   Position `json:"position" validate:"required"`
	Index      string   `json:"index" validate:"required"`
}

type UpdateStepInput struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"omitempty,min=2,max=255"`
	Timeout     int    `json:"timeout" validate:"omitempty,min=0"`
	EndpointID  string `json:"endpointId" validate:"omitempty,uuid"`
	WorkflowID  string `json:"workflowId" validate:"required,uuid"`
}

func NewStepOutput(step domain.Step) StepOutput {
	var endpoint *EndpointOutput
	if step.Endpoint.ID != uuid.Nil {
		endpointOutput := NewEndpointOutput(step.Endpoint)
		endpoint = &endpointOutput
	}

	return StepOutput{
		ID:          step.ID.String(),
		Name:        step.Name,
		Description: step.Description,
		Timeout:     step.Timeout,
		Position:    Position{X: step.Position.X, Y: step.Position.Y},
		Index:       step.Index,
		Endpoint:    endpoint,
		EndpointID:  step.EndpointID.String(),
		CreatedAt:   step.CreatedAt,
		UpdatedAt:   step.UpdatedAt,
	}
}

func NewStepsOutput(steps []domain.Step) []StepOutput {
	outputs := make([]StepOutput, len(steps))
	for i, step := range steps {
		outputs[i] = NewStepOutput(step)
	}
	return outputs
}
