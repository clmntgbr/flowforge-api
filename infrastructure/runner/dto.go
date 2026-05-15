package runner

import (
	"net/http"
	"time"
)

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

type ExecutionConfig struct {
	URL            string
	Method         string
	Headers        http.Header
	Body           []byte
	Timeout        int
	RetryOnFailure bool
	RetryCount     int
	RetryDelay     int
}

type RunnerResponse struct {
	Response string         `json:"response"`
	Insights RunnerInsights `json:"insights"`
}
