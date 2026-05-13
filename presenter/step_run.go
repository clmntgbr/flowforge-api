package presenter

import (
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"time"

	"github.com/google/uuid"
)

type StepRunListResponse struct {
	ID          uuid.UUID          `json:"id"`
	Status      enum.StepRunStatus `json:"status"`
	StartedAt   *time.Time         `json:"started_at"`
	CompletedAt *time.Time         `json:"completed_at"`
	FailedAt    *time.Time         `json:"failed_at"`
}

type StepRunDetailResponse struct {
	ID          uuid.UUID          `json:"id"`
	Status      enum.StepRunStatus `json:"status"`
	StartedAt   *time.Time         `json:"started_at"`
	CompletedAt *time.Time         `json:"completed_at"`
	FailedAt    *time.Time         `json:"failed_at"`
	Error       string             `json:"error"`
	Response    string             `json:"response"`
	Insight     InsightResponse    `json:"insight"`
	Step        StepDetailResponse `json:"step"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func NewStepRunListResponse(stepRun entity.StepRun) StepRunListResponse {
	return StepRunListResponse{
		ID:          stepRun.ID,
		Status:      stepRun.Status,
		StartedAt:   stepRun.StartedAt,
		CompletedAt: stepRun.CompletedAt,
		FailedAt:    stepRun.FailedAt,
	}
}

func NewStepRunListResponses(stepRuns []entity.StepRun) []StepRunListResponse {
	responses := make([]StepRunListResponse, len(stepRuns))
	for i, stepRun := range stepRuns {
		responses[i] = NewStepRunListResponse(stepRun)
	}
	return responses
}

func NewStepRunDetailResponses(stepRuns []entity.StepRun) []StepRunDetailResponse {
	responses := make([]StepRunDetailResponse, len(stepRuns))
	for i, stepRun := range stepRuns {
		responses[i] = NewStepRunDetailResponse(stepRun)
	}
	return responses
}

func NewStepRunDetailResponse(stepRun entity.StepRun) StepRunDetailResponse {
	return StepRunDetailResponse{
		ID:          stepRun.ID,
		Status:      stepRun.Status,
		StartedAt:   stepRun.StartedAt,
		CompletedAt: stepRun.CompletedAt,
		Error:       stepRun.Error,
		FailedAt:    stepRun.FailedAt,
		CreatedAt:   stepRun.CreatedAt,
		UpdatedAt:   stepRun.UpdatedAt,
		Insight:     NewInsightResponse(stepRun.Insight),
		Step:        NewStepDetailResponse(stepRun.Step),
		Response:    stepRun.Response,
	}
}
