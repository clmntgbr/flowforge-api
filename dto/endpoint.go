package dto

import (
	"forgeflow-api/domain"
	"time"
)

type EndpointOutput struct {
	MinimalEndpointOutput
	BaseURI   string    `json:"baseUri"`
	Timeout   int       `json:"timeout"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	Name    string `json:"name" validate:"required,min=2,max=255"`
	BaseURI string `json:"baseUri" validate:"required,url"`
	Path    string `json:"path" validate:"required"`
	Method  string `json:"method" validate:"required"`
	Timeout int    `json:"timeout" validate:"required,min=1,max=300000,number"`
}

type UpdateEndpointInput struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	BaseURI string `json:"baseUri" validate:"required,url"`
	Path    string `json:"path" validate:"required"`
	Method  string `json:"method" validate:"required"`
	Timeout int    `json:"timeout" validate:"required,min=1,max=300000,number"`
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
	}
}

func NewMinimalEndpointsOutput(endpoints []domain.Endpoint) []MinimalEndpointOutput {
	outputs := make([]MinimalEndpointOutput, len(endpoints))
	for i, endpoint := range endpoints {
		outputs[i] = NewMinimalEndpointOutput(endpoint)
	}
	return outputs
}
