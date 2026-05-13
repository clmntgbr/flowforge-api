package entity

import (
	"time"

	"github.com/google/uuid"
)

type Insight struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	StartTime time.Time     `gorm:"null" json:"start_time"`
	EndTime   time.Time     `gorm:"null" json:"end_time"`
	QueueTime time.Duration `gorm:"null" json:"queue_time"`

	DNSLookupDuration time.Duration `gorm:"null" json:"dns_lookup_duration"`
	TCPConnectionTime time.Duration `gorm:"null" json:"tcp_connection_time"`
	TLSHandshakeTime  time.Duration `gorm:"null" json:"tls_handshake_time"`
	TTFB              time.Duration `gorm:"null" json:"ttfb"`

	Duration      time.Duration `gorm:"null" json:"duration"`
	StatusCode    int           `gorm:"null;index:idx_insight_status" json:"status_code"`
	ResponseSize  int64         `gorm:"null" json:"response_size"`
	RequestSize   int64         `gorm:"null" json:"request_size"`
	AttemptNumber int           `gorm:"null" json:"attempt_number"`
	TotalAttempts int           `gorm:"nul" json:"total_attempts"`
	ErrorMessage  string        `gorm:"null" json:"error_message"`
	ErrorType     string        `gorm:"null;index:idx_insight_error_type" json:"error_type"`

	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_insight_created" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Insight) TableName() string {
	return "insights"
}
