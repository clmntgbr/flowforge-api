package workflow

import "flowforge-api/infrastructure/step"

type CreateWorkflowInput struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty,min=2,max=255"`
}

type UpdateWorkflowInput struct {
	CreateWorkflowInput
}

type UpsertWorkflowInput struct {
	Steps []step.UpsertWorkflowStepInput `json:"steps" validate:"omitempty,dive"`
}
