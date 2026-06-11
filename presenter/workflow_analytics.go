package presenter

import (
	"math"
	"time"
)

type WorkflowAnalyticsResponse struct {
	TotalRuns       int     `json:"totalRuns"`
	SuccessRate     float64 `json:"successRate"`
	SuccessCount    int     `json:"successCount"`
	FailureRate     float64 `json:"failureRate"`
	FailureCount    int     `json:"failureCount"`
	AverageDuration float64 `json:"averageDuration"`
}

func NewWorkflowAnalyticsResponse(totalRuns int, successRate float64, successCount int, failureRate float64, failureCount int, averageDuration time.Duration) WorkflowAnalyticsResponse {
	return WorkflowAnalyticsResponse{
		TotalRuns:       totalRuns,
		SuccessRate:     successRate,
		SuccessCount:    successCount,
		FailureRate:     failureRate,
		FailureCount:    failureCount,
		AverageDuration: math.Round(averageDuration.Seconds()*100) / 100,
	}
}
