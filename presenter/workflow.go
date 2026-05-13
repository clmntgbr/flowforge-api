package presenter

import (
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"fmt"
	"time"
)

type WorkflowListResponse struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Status    enum.WorkflowStatus `json:"status"`
	IsActive  bool                `json:"isActive"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

type WorkflowDetailResponse struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Status      enum.WorkflowStatus       `json:"status"`
	IsActive    bool                      `json:"isActive"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
	Description string                    `json:"description"`
	Steps       []StepDetailResponse      `json:"steps"`
	Connexions  []ConnexionDetailResponse `json:"connexions"`
}

func NewWorkflowListResponse(workflow entity.Workflow) WorkflowListResponse {
	return WorkflowListResponse{
		ID:        workflow.ID.String(),
		Name:      workflow.Name,
		CreatedAt: workflow.CreatedAt,
		UpdatedAt: workflow.UpdatedAt,
		Status:    workflow.Status,
		IsActive:  workflow.Status == enum.WorkflowStatusActive,
	}
}

func NewWorkflowListResponses(workflows []entity.Workflow) []WorkflowListResponse {
	responses := make([]WorkflowListResponse, len(workflows))
	for i, workflow := range workflows {
		responses[i] = NewWorkflowListResponse(workflow)
	}
	return responses
}

func NewWorkflowDetailResponse(workflow entity.Workflow) WorkflowDetailResponse {
	fmt.Println(workflow.Status)
	return WorkflowDetailResponse{
		ID:          workflow.ID.String(),
		Name:        workflow.Name,
		Status:      workflow.Status,
		CreatedAt:   workflow.CreatedAt,
		UpdatedAt:   workflow.UpdatedAt,
		Description: workflow.Description,
		Steps:       NewStepDetailResponses(workflow.Steps),
		Connexions:  NewConnexionDetailResponses(workflow.Connexions),
	}
}
