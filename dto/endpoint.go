package dto

import (
	"forgeflow-api/domain"
	"time"
)

type EndpointOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BaseURI   string    `json:"baseUri"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Timeout   int       `json:"timeout"`
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

func NewEndpointOutput(endpoint domain.Endpoint) EndpointOutput {
	return EndpointOutput{
		ID:        endpoint.ID.String(),
		Name:      endpoint.Name,
		BaseURI:   endpoint.BaseURI,
		Path:      endpoint.Path,
		Method:    endpoint.Method,
		Timeout:   endpoint.Timeout,
		CreatedAt: endpoint.CreatedAt,
		UpdatedAt: endpoint.UpdatedAt,
	}
}

func NewEndpointsOutput(endpoints []domain.Endpoint) []EndpointOutput {
	outputs := make([]EndpointOutput, len(endpoints))
	for i, endpoint := range endpoints {
		outputs[i] = NewEndpointOutput(endpoint)
	}
	return outputs
}
