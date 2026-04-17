package dto

import (
	"forgeflow-api/domain"

	"github.com/google/uuid"
)

type OrganizationOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type CreateOrganizationInput struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type UpdateOrganizationInput struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

func NewOrganizationOutput(organization domain.Organization, activeOrganizationID uuid.UUID) OrganizationOutput {
	isActive := false
	if activeOrganizationID != uuid.Nil && activeOrganizationID == organization.ID {
		isActive = true
	}

	return OrganizationOutput{
		ID:       organization.ID.String(),
		Name:     organization.Name,
		IsActive: isActive,
	}
}

func NewOrganizationsOutput(organizations []domain.Organization, activeOrganizationID uuid.UUID) []OrganizationOutput {
	outputs := make([]OrganizationOutput, len(organizations))
	for i, organization := range organizations {
		outputs[i] = NewOrganizationOutput(organization, activeOrganizationID)
	}
	return outputs
}
