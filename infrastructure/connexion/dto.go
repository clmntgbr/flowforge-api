package connexion

import "github.com/google/uuid"

type CreateConnexionInput struct {
	WorkflowID uuid.UUID `json:"workflowId" validate:"required,uuid"`
	FromStepID uuid.UUID `json:"fromStepId" validate:"required,uuid"`
	ToStepID   uuid.UUID `json:"toStepId" validate:"required,uuid"`
}
