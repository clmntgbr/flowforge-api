package presenter

import (
	"flowforge-api/domain/entity"
	"time"
)

type OrganizationListResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"isActive"`
}

type OrganizationDetailResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewOrganizationListResponse(org entity.Organization) OrganizationListResponse {
	return OrganizationListResponse{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
		IsActive:  org.IsActive,
	}
}

func NewOrganizationListResponses(orgs []entity.Organization) []OrganizationListResponse {
	responses := make([]OrganizationListResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = NewOrganizationListResponse(org)
	}
	return responses
}

func NewOrganizationDetailResponse(org entity.Organization) OrganizationDetailResponse {
	return OrganizationDetailResponse{
		ID:        org.ID.String(),
		Name:      org.Name,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}
}
