package presenter

import (
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"time"
)

type WorkflowRunResponse struct {
	ID          string                   `json:"id"`
	Status      enum.WorkflowRunStatus   `json:"status"`
	Statuses    []enum.WorkflowRunStatus `json:"statuses"`
	StartedAt   *time.Time               `json:"started_at"`
	CompletedAt *time.Time               `json:"completed_at"`
	StepsRuns   []StepRunDetailResponse  `json:"steps_runs"`
	FailedAt    *time.Time               `json:"failed_at"`
	CanceledAt  *time.Time               `json:"canceled_at"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

func NewWorkflowRunResponses(workflowRuns []entity.WorkflowRun) []WorkflowRunResponse {
	responses := make([]WorkflowRunResponse, len(workflowRuns))
	for i, workflowRun := range workflowRuns {
		responses[i] = NewWorkflowRunResponse(workflowRun)
	}
	return responses
}

func NewWorkflowRunResponse(workflowRun entity.WorkflowRun) WorkflowRunResponse {
	return WorkflowRunResponse{
		ID:          workflowRun.ID.String(),
		Status:      workflowRun.Status,
		Statuses:    workflowRun.Statuses,
		StartedAt:   workflowRun.StartedAt,
		CompletedAt: workflowRun.CompletedAt,
		FailedAt:    workflowRun.FailedAt,
		CanceledAt:  workflowRun.CanceledAt,
		CreatedAt:   workflowRun.CreatedAt,
		UpdatedAt:   workflowRun.UpdatedAt,
		StepsRuns:   NewStepRunDetailResponses(workflowRun.StepsRuns),
	}
}
