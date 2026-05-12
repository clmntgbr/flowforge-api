package presenter

import (
	"flowforge-api/domain/entity"
	"time"

	"github.com/google/uuid"
)

type OrganizationResponse struct {
	MinimalOrganizationResponse
}

type MinimalOrganizationResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewMinimalOrganizationsOutput(organizations []entity.Organization, activeOrganizationID uuid.UUID) []MinimalOrganizationResponse {
	outputs := make([]MinimalOrganizationResponse, len(organizations))
	for i, organization := range organizations {
		outputs[i] = NewMinimalOrganizationResponse(organization, activeOrganizationID)
	}
	return outputs
}

func NewMinimalOrganizationResponse(organization entity.Organization, activeOrganizationID uuid.UUID) MinimalOrganizationResponse {
	isActive := false
	if activeOrganizationID != uuid.Nil && activeOrganizationID == organization.ID {
		isActive = true
	}

	return MinimalOrganizationResponse{
		ID:        organization.ID.String(),
		Name:      organization.Name,
		IsActive:  isActive,
		CreatedAt: organization.CreatedAt,
		UpdatedAt: organization.UpdatedAt,
	}
}

func NewOrganizationResponse(organization entity.Organization, activeOrganizationID uuid.UUID) OrganizationResponse {
	return OrganizationResponse{
		MinimalOrganizationResponse: NewMinimalOrganizationResponse(organization, activeOrganizationID),
	}
}
