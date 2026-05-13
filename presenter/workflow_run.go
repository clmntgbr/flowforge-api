package presenter

import (
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"time"
)

type WorkflowRunListResponse struct {
	ID          string                 `json:"id"`
	Status      enum.WorkflowRunStatus `json:"status"`
	StartedAt   *time.Time             `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	Error       string                 `json:"error"`
	StepsRuns   []StepRunListResponse  `json:"steps_runs"`
	FailedAt    *time.Time             `json:"failed_at"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type WorkflowRunDetailResponse struct {
	ID          string                  `json:"id"`
	Status      enum.WorkflowRunStatus  `json:"status"`
	StartedAt   *time.Time              `json:"started_at"`
	CompletedAt *time.Time              `json:"completed_at"`
	Error       string                  `json:"error"`
	StepsRuns   []StepRunDetailResponse `json:"steps_runs"`
	FailedAt    *time.Time              `json:"failed_at"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

func NewWorkflowRunListResponse(workflowRun entity.WorkflowRun) WorkflowRunListResponse {
	return WorkflowRunListResponse{
		ID:          workflowRun.ID.String(),
		Status:      workflowRun.Status,
		StartedAt:   workflowRun.StartedAt,
		CompletedAt: workflowRun.CompletedAt,
		Error:       workflowRun.Error,
		FailedAt:    workflowRun.FailedAt,
		CreatedAt:   workflowRun.CreatedAt,
		UpdatedAt:   workflowRun.UpdatedAt,
		StepsRuns:   NewStepRunListResponses(workflowRun.StepsRuns),
	}
}

func NewWorkflowRunListResponses(workflowRuns []entity.WorkflowRun) []WorkflowRunListResponse {
	responses := make([]WorkflowRunListResponse, len(workflowRuns))
	for i, workflowRun := range workflowRuns {
		responses[i] = NewWorkflowRunListResponse(workflowRun)
	}
	return responses
}

func NewWorkflowRunDetailResponse(workflowRun entity.WorkflowRun) WorkflowRunDetailResponse {
	return WorkflowRunDetailResponse{
		ID:          workflowRun.ID.String(),
		Status:      workflowRun.Status,
		StartedAt:   workflowRun.StartedAt,
		CompletedAt: workflowRun.CompletedAt,
		Error:       workflowRun.Error,
		FailedAt:    workflowRun.FailedAt,
		CreatedAt:   workflowRun.CreatedAt,
		UpdatedAt:   workflowRun.UpdatedAt,
		StepsRuns:   NewStepRunDetailResponses(workflowRun.StepsRuns),
	}
}
