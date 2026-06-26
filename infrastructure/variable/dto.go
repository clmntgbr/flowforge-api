package workflow

import "github.com/google/uuid"

type CreateVariableInput struct {
	Name        string    `json:"name" validate:"required,min=2,max=255"`
	Description string    `json:"description" validate:"omitempty,min=2,max=255"`
	Path        string    `json:"path" validate:"required,min=2,max=255"`
	StepID      uuid.UUID `json:"stepId" validate:"required,uuid"`
}

type UpdateVariableInput struct {
	CreateVariableInput
}
