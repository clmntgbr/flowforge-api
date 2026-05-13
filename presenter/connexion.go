package presenter

import (
	"flowforge-api/domain/entity"
	"time"
)

type ConnexionDetailResponse struct {
	ID         string    `json:"id"`
	FromStepID string    `json:"fromStepId"`
	ToStepID   string    `json:"toStepId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewConnexionDetailResponse(connexion entity.Connexion) ConnexionDetailResponse {
	return ConnexionDetailResponse{
		ID:         connexion.ID.String(),
		FromStepID: connexion.FromStepID.String(),
		ToStepID:   connexion.ToStepID.String(),
		CreatedAt:  connexion.CreatedAt,
		UpdatedAt:  connexion.UpdatedAt,
	}
}

func NewConnexionDetailResponses(connexions []entity.Connexion) []ConnexionDetailResponse {
	responses := make([]ConnexionDetailResponse, len(connexions))
	for i, connexion := range connexions {
		responses[i] = NewConnexionDetailResponse(connexion)
	}
	return responses
}
