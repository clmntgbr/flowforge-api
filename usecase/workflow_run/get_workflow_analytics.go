package workflow_run

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/enum"
	"flowforge-api/domain/repository"
	"math"
	"time"

	"github.com/samber/lo"

	"github.com/google/uuid"
)

type GetWorkflowAnalyticsUseCase struct {
	workflowRepo    *repository.WorkflowRepository
	workflowRunRepo *repository.WorkflowRunRepository
}

func NewGetWorkflowAnalyticsUseCase(
	workflowRepo *repository.WorkflowRepository,
	workflowRunRepo *repository.WorkflowRunRepository,
) *GetWorkflowAnalyticsUseCase {
	return &GetWorkflowAnalyticsUseCase{
		workflowRepo:    workflowRepo,
		workflowRunRepo: workflowRunRepo,
	}
}

func (u *GetWorkflowAnalyticsUseCase) Execute(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (int, float64, int, float64, int, time.Duration, error) {
	workflow, err := (*u.workflowRepo).GetByIDAndOrganizationID(ctx, organizationID, workflowID)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	workflowRuns, err := (*u.workflowRunRepo).GetAllByWorkflowID(ctx, workflow.ID)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	totalRuns := calculateTotalRuns(workflowRuns)

	successCount := calculateSuccessCount(workflowRuns)
	failureCount := calculateFailureCount(workflowRuns)

	successRate := calculateSuccessRate(successCount, totalRuns)
	failureRate := calculateFailureRate(failureCount, totalRuns)

	averageDuration := calculateAverageDuration(workflowRuns)

	return totalRuns, successRate, successCount, failureRate, failureCount, averageDuration, nil
}

func calculateSuccessRate(successCount int, totalCount int) float64 {
	if totalCount == 0 {
		return 0
	}
	return math.Round(float64(successCount) / float64(totalCount) * 100)
}

func calculateFailureRate(failureCount int, totalCount int) float64 {
	if totalCount == 0 {
		return 0
	}
	return math.Round(float64(failureCount) / float64(totalCount) * 100)
}

func calculateAverageDuration(workflowRuns []entity.WorkflowRun) time.Duration {
	if len(workflowRuns) == 0 {
		return 0
	}
	totalDuration := time.Duration(0)
	counted := 0
	for _, run := range workflowRuns {
		if run.StartedAt == nil {
			continue
		}
		if run.CompletedAt == nil {
			continue
		}
		totalDuration += run.CompletedAt.Sub(*run.StartedAt)
		counted++
	}
	if counted == 0 {
		return 0
	}
	avg := totalDuration / time.Duration(counted)
	return avg.Round(10 * time.Millisecond)
}

func calculateTotalRuns(workflowRuns []entity.WorkflowRun) int {
	return len(workflowRuns)
}

func calculateSuccessCount(workflowRuns []entity.WorkflowRun) int {
	return lo.CountBy(workflowRuns, func(run entity.WorkflowRun) bool {
		return run.Status == enum.WorkflowRunStatusCompleted
	})
}

func calculateFailureCount(workflowRuns []entity.WorkflowRun) int {
	return lo.CountBy(workflowRuns, func(run entity.WorkflowRun) bool {
		return run.Status == enum.WorkflowRunStatusFailed
	})
}
