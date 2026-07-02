package presenter

import (
	"flowforge-api/domain/entity"
)

type VariableResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Path         string             `json:"path"`
	StepID       string             `json:"stepId"`
	Step         StepDetailResponse `json:"step"`
	IsSecret     bool               `gorm:"default:false" json:"isSecret"`
	DefaultValue string             `gorm:"null" json:"defaultValue"`
	LastValue    string             `gorm:"null" json:"lastValue"`
}

func NewVariableResponses(variables []entity.Variable) []VariableResponse {
	responses := make([]VariableResponse, len(variables))
	for i, variable := range variables {
		responses[i] = NewVariableResponse(variable)
	}
	return responses
}

func NewVariableResponse(variable entity.Variable) VariableResponse {
	return VariableResponse{
		ID:           variable.ID.String(),
		Name:         variable.Name,
		Path:         variable.Path,
		StepID:       variable.StepID.String(),
		Step:         NewStepDetailResponse(variable.Step),
		IsSecret:     variable.IsSecret,
		DefaultValue: variable.DefaultValue,
		LastValue:    variable.LastValue,
	}
}
