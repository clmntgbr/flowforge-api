package insight

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"time"
)

type CreateInsightUseCase struct {
	insightRepository *repository.InsightRepository
}

func NewCreateInsightUseCase(insightRepository *repository.InsightRepository) *CreateInsightUseCase {
	return &CreateInsightUseCase{
		insightRepository: insightRepository,
	}
}

func (u *CreateInsightUseCase) Execute(
	ctx context.Context,
	startTime time.Time,
	endTime time.Time,
	duration time.Duration,
	statusCode int,
	responseSize int64,
	attemptNumber int,
	totalAttempts int,
	queueTime time.Duration,
	dnsLookupDuration time.Duration,
	tcpConnectionTime time.Duration,
	tlsHandshakeTime time.Duration,
	ttfb time.Duration,
	errorMessage string,
	errorType string,
	requestSize int64,
) (*entity.Insight, error) {
	insight := &entity.Insight{
		StartTime:         startTime,
		EndTime:           endTime,
		Duration:          duration,
		StatusCode:        statusCode,
		ResponseSize:      responseSize,
		AttemptNumber:     attemptNumber,
		TotalAttempts:     totalAttempts,
		QueueTime:         queueTime,
		DNSLookupDuration: dnsLookupDuration,
		TCPConnectionTime: tcpConnectionTime,
		TLSHandshakeTime:  tlsHandshakeTime,
		TTFB:              ttfb,
		ErrorMessage:      errorMessage,
		ErrorType:         errorType,
		RequestSize:       requestSize,
	}

	err := (*u.insightRepository).Create(ctx, insight)
	if err != nil {
		return nil, err
	}

	return insight, nil
}
