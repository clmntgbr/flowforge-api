package workflow

import "github.com/google/uuid"

type CreateVariableInput struct {
	Name        string    `json:"name" validate:"required,min=2,max=255"`
	Key         string    `json:"key" validate:"required,min=2,max=255"`
	Description string    `json:"description" validate:"omitempty,min=2,max=255"`
	Path        string    `json:"path" validate:"required,min=2,max=255"`
	StepID      uuid.UUID `json:"stepId" validate:"required,uuid"`
}

type UpdateVariableInput struct {
	CreateVariableInput
}

type SearchVariablesPathInput struct {
	StepID uuid.UUID `json:"stepId" validate:"required,uuid"`
	Query  string    `json:"query" validate:"omitempty,min=1,max=255"`
	Page   int       `json:"page" query:"page"`
	Limit  int       `json:"limit" query:"limit"`
}

func (s *SearchVariablesPathInput) Normalize() {
	if s.Page <= 0 {
		s.Page = 1
	}
	if s.Limit <= 0 || s.Limit > 100 {
		s.Limit = 20
	}
}

func (s *SearchVariablesPathInput) Offset() int {
	return (s.Page - 1) * s.Limit
}
