package dto

import (
	"forgeflow-api/domain"
	"time"

	"github.com/google/uuid"
)

type ConnexionOutput struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateConnexionInput struct {
	WorkflowID uuid.UUID `json:"workflowId" validate:"required,uuid"`
	From       string    `json:"from" validate:"required,uuid"`
	To         string    `json:"to" validate:"required,uuid"`
}

type DeleteConnexionInput struct {
	WorkflowID string `json:"workflowId" validate:"required,uuid"`
	ID         string `json:"id" validate:"required,uuid"`
}

func NewConnexionOutput(connexion domain.Connexion) ConnexionOutput {
	return ConnexionOutput{
		ID:        connexion.ID.String(),
		From:      connexion.FromStepID.String(),
		To:        connexion.ToStepID.String(),
		CreatedAt: connexion.CreatedAt,
		UpdatedAt: connexion.UpdatedAt,
	}
}

func NewConnexionsOutput(connexions []domain.Connexion) []ConnexionOutput {
	outputs := make([]ConnexionOutput, len(connexions))
	for i, connexion := range connexions {
		outputs[i] = NewConnexionOutput(connexion)
	}
	return outputs
}
