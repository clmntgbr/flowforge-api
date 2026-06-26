package presenter

import (
	"flowforge-api/domain/entity"
)

type VariableResponse struct {
	ID     string             `json:"id"`
	Name   string             `json:"name"`
	Path   string             `json:"path"`
	StepID string             `json:"stepId"`
	Step   StepDetailResponse `json:"step"`
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
		ID:     variable.ID.String(),
		Name:   variable.Name,
		Path:   variable.Path,
		StepID: variable.StepID.String(),
		Step:   NewStepDetailResponse(variable.Step),
	}
}
