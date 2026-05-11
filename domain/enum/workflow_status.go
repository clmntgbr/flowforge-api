package enum

type WorkflowStatus string

const (
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusInactive WorkflowStatus = "inactive"
	WorkflowStatusDeleted  WorkflowStatus = "deleted"
)
