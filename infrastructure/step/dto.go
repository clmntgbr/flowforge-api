package step

import "flowforge-api/domain/types"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type UpdateStepInput struct {
	Name           string       `json:"name" validate:"required,min=2,max=100"`
	Description    string       `json:"description" validate:"omitempty,min=2,max=255"`
	Timeout        int          `json:"timeout" validate:"required,min=1,max=60,number"`
	Query          types.Query  `json:"query" validate:"required,dive"`
	Header         types.Header `json:"header" validate:"required,dive"`
	Body           types.Body   `json:"body" validate:"required,dive"`
	RetryOnFailure bool         `json:"retryOnFailure"`
	RetryCount     int          `json:"retryCount" validate:"min=0,max=10,number"`
	RetryDelay     int          `json:"retryDelay" validate:"min=0,max=600,number"`
}

type UpsertWorkflowStepInput struct {
	ID         string   `json:"id" validate:"required,uuid"`
	EndpointID string   `json:"endpointId" validate:"required,uuid"`
	Position   Position `json:"position" validate:"required"`
	Index      string   `json:"index" validate:"required"`
}
