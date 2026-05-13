package presenter

import (
	"flowforge-api/domain/entity"
	"time"

	"github.com/google/uuid"
)

type InsightResponse struct {
	ID uuid.UUID `json:"id"`

	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	QueueTime time.Duration `json:"queue_time"`

	DNSLookupDuration time.Duration `json:"dns_lookup_duration"`
	TCPConnectionTime time.Duration `json:"tcp_connection_time"`
	TLSHandshakeTime  time.Duration `json:"tls_handshake_time"`
	TTFB              time.Duration `json:"ttfb"`

	Duration      time.Duration `json:"duration"`
	StatusCode    int           `json:"status_code"`
	ResponseSize  int64         `json:"response_size"`
	RequestSize   int64         `json:"request_size"`
	AttemptNumber int           `json:"attempt_number"`
	TotalAttempts int           `json:"total_attempts"`
	ErrorMessage  string        `json:"error_message"`
	ErrorType     string        `json:"error_type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewInsightResponse(insight *entity.Insight) InsightResponse {
	if insight == nil {
		return InsightResponse{}
	}

	return InsightResponse{
		ID:                insight.ID,
		StartTime:         insight.StartTime,
		EndTime:           insight.EndTime,
		QueueTime:         insight.QueueTime,
		DNSLookupDuration: insight.DNSLookupDuration,
		TCPConnectionTime: insight.TCPConnectionTime,
		TLSHandshakeTime:  insight.TLSHandshakeTime,
		TTFB:              insight.TTFB,
		Duration:          insight.Duration,
		StatusCode:        insight.StatusCode,
		ResponseSize:      insight.ResponseSize,
		RequestSize:       insight.RequestSize,
		AttemptNumber:     insight.AttemptNumber,
		TotalAttempts:     insight.TotalAttempts,
		ErrorMessage:      insight.ErrorMessage,
		ErrorType:         insight.ErrorType,
		CreatedAt:         insight.CreatedAt,
		UpdatedAt:         insight.UpdatedAt,
	}
}

func NewInsightResponses(insights []entity.Insight) []InsightResponse {
	responses := make([]InsightResponse, len(insights))
	for i, insight := range insights {
		responses[i] = NewInsightResponse(&insight)
	}
	return responses
}
