package enum

type StepRunStatus string

const (
	StepRunStatusPending   StepRunStatus = "pending"
	StepRunStatusRunning   StepRunStatus = "running"
	StepRunStatusCompleted StepRunStatus = "completed"
	StepRunStatusFailed    StepRunStatus = "failed"
	StepRunStatusCanceled  StepRunStatus = "canceled"
)
