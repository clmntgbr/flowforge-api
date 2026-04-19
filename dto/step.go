package dto

import (
	"forgeflow-api/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type StepOutput struct {
	MinimalStepOutput
	Description    string         `json:"description"`
	Timeout        int            `json:"timeout"`
	Query          datatypes.JSON `json:"query"`
	RetryOnFailure bool           `json:"retryOnFailure"`
	RetryCount     int            `json:"retryCount"`
	RetryDelay     int            `json:"retryDelay"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type MinimalStepOutput struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Position   Position        `json:"position"`
	Index      string          `json:"index"`
	Endpoint   *EndpointOutput `json:"endpoint,omitempty"`
	EndpointID string          `json:"endpointId"`
}

type UpdateWorkflowStepInput struct {
	ID         string   `json:"id" validate:"required,uuid"`
	EndpointID string   `json:"endpointId" validate:"required,uuid"`
	Position   Position `json:"position" validate:"required"`
	Index      string   `json:"index" validate:"required"`
}

type UpdateStepInput struct {
	Name           string         `json:"name" validate:"required,min=2,max=100"`
	Description    string         `json:"description" validate:"omitempty,min=2,max=255"`
	Timeout        int            `json:"timeout" validate:"omitempty,min=0"`
	Query          datatypes.JSON `json:"query"`
	RetryOnFailure bool           `json:"retryOnFailure"`
	RetryCount     int            `json:"retryCount" validate:"min=0,max=10,number"`
	RetryDelay     int            `json:"retryDelay" validate:"min=0,max=300000,number"`
}

func NewStepOutput(step domain.Step) StepOutput {
	return StepOutput{
		MinimalStepOutput: NewMinimalStepOutput(step),
		Description:       step.Description,
		Timeout:           step.Timeout,
		Query:             step.Query,
		RetryOnFailure:    step.RetryOnFailure,
		RetryCount:        step.RetryCount,
		RetryDelay:        step.RetryDelay,
		CreatedAt:         step.CreatedAt,
		UpdatedAt:         step.UpdatedAt,
	}
}

func NewMinimalStepOutput(step domain.Step) MinimalStepOutput {
	var endpoint *EndpointOutput
	if step.Endpoint.ID != uuid.Nil {
		endpointOutput := NewEndpointOutput(step.Endpoint)
		endpoint = &endpointOutput
	}

	return MinimalStepOutput{
		ID:         step.ID.String(),
		Name:       step.Name,
		Position:   Position{X: step.Position.X, Y: step.Position.Y},
		Index:      step.Index,
		Endpoint:   endpoint,
		EndpointID: step.EndpointID.String(),
	}
}

func NewMinimalStepsOutput(steps []domain.Step) []MinimalStepOutput {
	outputs := make([]MinimalStepOutput, len(steps))
	for i, step := range steps {
		outputs[i] = NewMinimalStepOutput(step)
	}
	return outputs
}

func NewStepsOutput(steps []domain.Step) []StepOutput {
	outputs := make([]StepOutput, len(steps))
	for i, step := range steps {
		outputs[i] = NewStepOutput(step)
	}
	return outputs
}
