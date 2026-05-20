package endpoint

import "flowforge-api/domain/types"

type CreateEndpointInput struct {
	Name           string       `json:"name" validate:"required,min=2,max=255"`
	BaseURI        string       `json:"baseUri" validate:"required,url"`
	Path           string       `json:"path" validate:"required"`
	Method         string       `json:"method" validate:"required"`
	Timeout        int          `json:"timeout" validate:"required,min=1,max=60,number"`
	Query          types.Query  `json:"query" validate:"required,dive"`
	Header         types.Header `json:"header" validate:"required,dive"`
	Body           types.Body   `json:"body" validate:"required,dive"`
	RetryOnFailure bool         `json:"retryOnFailure"`
	RetryCount     int          `json:"retryCount" validate:"min=0,max=10,number"`
	RetryDelay     int          `json:"retryDelay" validate:"min=0,max=600,number"`
}

type OpenAPI struct {
	OpenAPI string `json:"openapi"`
	Info    struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Servers []struct {
		URL string `json:"url"`
	} `json:"servers"`
	Paths map[string]map[string]struct {
		OperationID string   `json:"operationId"`
		Tags        []string `json:"tags"`
		Summary     string   `json:"summary"`
	} `json:"paths"`
}

type UpdateEndpointInput struct {
	CreateEndpointInput
}
