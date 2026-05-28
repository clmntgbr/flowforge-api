package entity

import (
	"flowforge-api/domain/types"
	"time"

	"github.com/google/uuid"
)

type Insight struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	StartTime time.Time      `gorm:"null" json:"start_time"`
	EndTime   time.Time      `gorm:"null" json:"end_time"`
	QueueTime types.Duration `gorm:"type:bigint;null" json:"queue_time"`

	DNSLookupDuration types.Duration `gorm:"type:bigint;null" json:"dns_lookup_duration"`
	TCPConnectionTime types.Duration `gorm:"type:bigint;null" json:"tcp_connection_time"`
	TLSHandshakeTime  types.Duration `gorm:"type:bigint;null" json:"tls_handshake_time"`
	TTFB              types.Duration `gorm:"type:bigint;null" json:"ttfb"`

	Duration types.Duration `gorm:"type:bigint;null" json:"duration"`
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
