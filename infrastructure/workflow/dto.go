package workflow

import "flowforge-api/infrastructure/step"

type CreateWorkflowInput struct {
	Name                 string `json:"name" validate:"required,min=2,max=255"`
	Description          string `json:"description" validate:"omitempty,min=2,max=255"`
	NotificationsEnabled bool   `json:"notificationsEnabled" validate:"omitempty"`
	NotifyOnSuccess      bool   `json:"notifyOnSuccess" validate:"omitempty"`
	NotifyOnFailure      bool   `json:"notifyOnFailure" validate:"omitempty"`
	NotifyOnCancel       bool   `json:"notifyOnCancel" validate:"omitempty"`
}

type UpdateWorkflowInput struct {
	CreateWorkflowInput
}

type UpsertWorkflowInput struct {
	Steps []step.UpsertWorkflowStepInput `json:"steps" validate:"omitempty,dive"`
}
