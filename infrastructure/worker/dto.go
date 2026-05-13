package worker

import "time"

type RunnerInsights struct {
	StartTime         time.Time
	EndTime           time.Time
	Duration          time.Duration
	StatusCode        int
	ResponseSize      int64
	AttemptNumber     int
	TotalAttempts     int
	QueueTime         time.Duration
	DNSLookupDuration time.Duration
	TCPConnectionTime time.Duration
	TLSHandshakeTime  time.Duration
	TTFB              time.Duration
	ErrorMessage      string
	ErrorType         string
	RequestSize       int64
}
