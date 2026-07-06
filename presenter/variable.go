package presenter

import (
	"flowforge-api/domain/entity"
)

type VariableListResponse struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Key    string           `json:"key"`
	Path   string           `json:"path"`
	StepID string           `json:"stepId"`
	Step   StepListResponse `json:"step"`
}

type VariableResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Key          string             `json:"key"`
	Path         string             `json:"path"`
	StepID       string             `json:"stepId"`
	Step         StepDetailResponse `json:"step"`
	IsSecret     bool               `gorm:"default:false" json:"isSecret"`
	DefaultValue string             `gorm:"null" json:"defaultValue"`
	LastValue    string             `gorm:"null" json:"lastValue"`
}

func NewVariableListResponses(variables []entity.Variable) []VariableListResponse {
	responses := make([]VariableListResponse, len(variables))
	for i, variable := range variables {
		responses[i] = NewVariableListResponse(variable)
	}
	return responses
}

func NewVariableListResponse(variable entity.Variable) VariableListResponse {
	return VariableListResponse{
		ID:     variable.ID.String(),
		Name:   variable.Name,
		Key:    variable.Key,
		Path:   variable.Path,
		StepID: variable.StepID.String(),
		Step:   NewStepListResponse(variable.Step),
	}
}

func NewVariableDetailResponses(variable entity.Variable) VariableResponse {
	return NewVariableDetailResponse(variable)
}

func NewVariableDetailResponse(variable entity.Variable) VariableResponse {
	return VariableResponse{
		ID:           variable.ID.String(),
		Name:         variable.Name,
		Key:          variable.Key,
		Path:         variable.Path,
		StepID:       variable.StepID.String(),
		Step:         NewStepDetailResponse(variable.Step),
		IsSecret:     variable.IsSecret,
		DefaultValue: variable.DefaultValue,
		LastValue:    variable.LastValue,
	}
}
