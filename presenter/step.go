package presenter

import (
	"flowforge-api/domain/entity"
	"flowforge-api/domain/types"
	"time"
)

type StepDetailResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Position       entity.Position        `json:"position"`
	Index          string                 `json:"index"`
	Endpoint       EndpointDetailResponse `json:"endpoint,omitempty"`
	EndpointID     string                 `json:"endpointId"`
	Description    string                 `json:"description"`
	Timeout        int                    `json:"timeout"`
	Query          types.Query            `json:"query"`
	Header         types.Header           `json:"header"`
	Body           types.Body             `json:"body"`
	RetryOnFailure bool                   `json:"retryOnFailure"`
	RetryCount     int                    `json:"retryCount"`
	RetryDelay     int                    `json:"retryDelay"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}

func NewStepDetailResponse(step entity.Step) StepDetailResponse {
	return StepDetailResponse{
		ID:             step.ID.String(),
		Name:           step.Name,
		Position:       step.Position,
		Index:          step.Index,
		Endpoint:       NewEndpointDetailResponse(step.Endpoint),
		EndpointID:     step.EndpointID.String(),
		Description:    step.Description,
		Timeout:        step.Timeout,
		Query:          step.Query,
		Header:         step.Header,
		Body:           step.Body,
		RetryOnFailure: step.RetryOnFailure,
		RetryCount:     step.RetryCount,
		RetryDelay:     step.RetryDelay,
		CreatedAt:      step.CreatedAt,
		UpdatedAt:      step.UpdatedAt,
	}
}

func NewStepDetailResponses(steps []entity.Step) []StepDetailResponse {
	responses := make([]StepDetailResponse, len(steps))
	for i, step := range steps {
		responses[i] = NewStepDetailResponse(step)
	}
	return responses
}
