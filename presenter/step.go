package presenter

import (
	"flowforge-api/domain/entity"
	"flowforge-api/domain/types"
	"time"
)

type StepListResponse struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	Position       entity.Position      `json:"position"`
	Index          string               `json:"index"`
	ExecutionOrder int                  `json:"order"`
	TreeIndex      int                  `json:"treeIndex"`
	Endpoint       EndpointListResponse `json:"endpoint"`
}

type StepDetailResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Position       entity.Position        `json:"position"`
	Index          string                 `json:"index"`
	ExecutionOrder int                    `json:"order"`
	TreeIndex      int                    `json:"treeIndex"`
	Endpoint       EndpointDetailResponse `json:"endpoint,omitempty"`
	EndpointID     string                 `json:"endpointId"`
	Description    string                 `json:"description"`
	Timeout        int                    `json:"timeout"`
	URL            string                 `json:"url"`
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
		ExecutionOrder: step.ExecutionOrder,
		TreeIndex:      step.TreeIndex,
		Endpoint:       NewEndpointDetailResponse(step.Endpoint),
		EndpointID:     step.EndpointID.String(),
		Description:    step.Description,
		Timeout:        step.Timeout,
		URL:            step.URL,
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

func NewStepListResponse(step entity.Step) StepListResponse {
	return StepListResponse{
		ID:             step.ID.String(),
		Name:           step.Name,
		Position:       step.Position,
		Index:          step.Index,
		ExecutionOrder: step.ExecutionOrder,
		TreeIndex:      step.TreeIndex,
		Endpoint:       NewEndpointListResponse(step.Endpoint),
	}
}

func NewStepDetailResponses(steps []entity.Step) []StepDetailResponse {
	responses := make([]StepDetailResponse, len(steps))
	for i, step := range steps {
		responses[i] = NewStepDetailResponse(step)
	}
	return responses
}
