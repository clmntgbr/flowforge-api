package presenter

import (
	"flowforge-api/domain/entity"
	"time"

	"github.com/google/uuid"
)

type OrganizationOutput struct {
	MinimalOrganizationOutput
}

type MinimalOrganizationOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewMinimalOrganizationsOutput(organizations []entity.Organization, activeOrganizationID uuid.UUID) []MinimalOrganizationOutput {
	outputs := make([]MinimalOrganizationOutput, len(organizations))
	for i, organization := range organizations {
		outputs[i] = NewMinimalOrganizationOutput(organization, activeOrganizationID)
	}
	return outputs
}

func NewMinimalOrganizationOutput(organization entity.Organization, activeOrganizationID uuid.UUID) MinimalOrganizationOutput {
	isActive := false
	if activeOrganizationID != uuid.Nil && activeOrganizationID == organization.ID {
		isActive = true
	}

	return MinimalOrganizationOutput{
		ID:        organization.ID.String(),
		Name:      organization.Name,
		IsActive:  isActive,
		CreatedAt: organization.CreatedAt,
		UpdatedAt: organization.UpdatedAt,
	}
}

func NewOrganizationOutput(organization entity.Organization, activeOrganizationID uuid.UUID) OrganizationOutput {
	return OrganizationOutput{
		MinimalOrganizationOutput: NewMinimalOrganizationOutput(organization, activeOrganizationID),
	}
}
