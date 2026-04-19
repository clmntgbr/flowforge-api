package dto

import (
	"forgeflow-api/domain"
	"time"

	"gorm.io/datatypes"
)

type EndpointOutput struct {
	MinimalEndpointOutput
	BaseURI        string         `json:"baseUri"`
	Timeout        int            `json:"timeout"`
	Query          datatypes.JSON `json:"query"`
	RetryOnFailure bool           `json:"retryOnFailure"`
	RetryCount     int            `json:"retryCount"`
	RetryDelay     int            `json:"retryDelay"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type MinimalEndpointOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateEndpointInput struct {
	Name           string         `json:"name" validate:"required,min=2,max=255"`
	BaseURI        string         `json:"baseUri" validate:"required,url"`
	Path           string         `json:"path" validate:"required"`
	Method         string         `json:"method" validate:"required"`
	Timeout        int            `json:"timeout" validate:"required,min=1,max=300000,number"`
	Query          datatypes.JSON `json:"query" validate:"required,json"`
	RetryOnFailure bool           `json:"retryOnFailure"`
	RetryCount     int            `json:"retryCount" validate:"required,min=0,max=10,number"`
	RetryDelay     int            `json:"retryDelay" validate:"required,min=0,max=300000,number"`
}

type UpdateEndpointInput struct {
	Name           string         `json:"name" validate:"required,min=2,max=255"`
	BaseURI        string         `json:"baseUri" validate:"required,url"`
	Path           string         `json:"path" validate:"required"`
	Method         string         `json:"method" validate:"required"`
	Timeout        int            `json:"timeout" validate:"required,min=1,max=300000,number"`
	Query          datatypes.JSON `json:"query" validate:"required,json"`
	RetryOnFailure bool           `json:"retryOnFailure"`
	RetryCount     int            `json:"retryCount" validate:"required,min=0,max=10,number"`
	RetryDelay     int            `json:"retryDelay" validate:"required,min=0,max=300000,number"`
}

func NewMinimalEndpointOutput(endpoint domain.Endpoint) MinimalEndpointOutput {
	return MinimalEndpointOutput{
		ID:        endpoint.ID.String(),
		Name:      endpoint.Name,
		Path:      endpoint.Path,
		Method:    endpoint.Method,
		CreatedAt: endpoint.CreatedAt,
		UpdatedAt: endpoint.UpdatedAt,
	}
}

func NewEndpointOutput(endpoint domain.Endpoint) EndpointOutput {
	return EndpointOutput{
		MinimalEndpointOutput: NewMinimalEndpointOutput(endpoint),
		BaseURI:               endpoint.BaseURI,
		Timeout:               endpoint.Timeout,
		Query:                 endpoint.Query,
		RetryOnFailure:        endpoint.RetryOnFailure,
		RetryCount:            endpoint.RetryCount,
		RetryDelay:            endpoint.RetryDelay,
	}
}

func NewMinimalEndpointsOutput(endpoints []domain.Endpoint) []MinimalEndpointOutput {
	outputs := make([]MinimalEndpointOutput, len(endpoints))
	for i, endpoint := range endpoints {
		outputs[i] = NewMinimalEndpointOutput(endpoint)
	}
	return outputs
}
