package presenter

import (
	"flowforge-api/domain/entity"
	"flowforge-api/domain/types"
	"time"
)

type EndpointListResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	BaseURI string `json:"baseUri"`
	Path    string `json:"path"`
	Method  string `json:"method"`
}

type EndpointDetailResponse struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	BaseURI        string       `json:"baseUri"`
	Path           string       `json:"path"`
	Method         string       `json:"method"`
	Timeout        int          `json:"timeout"`
	RetryOnFailure bool         `json:"retryOnFailure"`
	RetryCount     int          `json:"retryCount"`
	RetryDelay     int          `json:"retryDelay"`
	Query          types.Query  `json:"query"`
	Header         types.Header `json:"header"`
	Body           types.Body   `json:"body"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}

func NewEndpointListResponse(endpoint entity.Endpoint) EndpointListResponse {
	return EndpointListResponse{
		ID:      endpoint.ID.String(),
		Name:    endpoint.Name,
		BaseURI: endpoint.BaseURI,
		Path:    endpoint.Path,
		Method:  endpoint.Method,
	}
}

func NewEndpointListResponses(endpoints []entity.Endpoint) []EndpointListResponse {
	responses := make([]EndpointListResponse, len(endpoints))
	for i, endpoint := range endpoints {
		responses[i] = NewEndpointListResponse(endpoint)
	}
	return responses
}

func NewEndpointDetailResponse(endpoint entity.Endpoint) EndpointDetailResponse {
	return EndpointDetailResponse{
		ID:             endpoint.ID.String(),
		Name:           endpoint.Name,
		BaseURI:        endpoint.BaseURI,
		Path:           endpoint.Path,
		Method:         endpoint.Method,
		Timeout:        endpoint.Timeout,
		RetryOnFailure: endpoint.RetryOnFailure,
		RetryCount:     endpoint.RetryCount,
		RetryDelay:     endpoint.RetryDelay,
		Query:          endpoint.Query,
		Header:         endpoint.Header,
		Body:           endpoint.Body,
		CreatedAt:      endpoint.CreatedAt,
		UpdatedAt:      endpoint.UpdatedAt,
	}
}
