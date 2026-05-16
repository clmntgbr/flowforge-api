package connexion

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/connexion"

	"github.com/google/uuid"
)

type CreateConnexionUseCase struct {
	connexionRepo *repository.ConnexionRepository
}

func NewCreateConnexionUseCase(connexionRepo *repository.ConnexionRepository) *CreateConnexionUseCase {
	return &CreateConnexionUseCase{connexionRepo: connexionRepo}
}

func (u *CreateConnexionUseCase) Execute(ctx context.Context, organizationID uuid.UUID, input connexion.CreateConnexionInput) (entity.Connexion, error) {
	connexions, err := (*u.connexionRepo).GetByFromStepIDAndToStepIDAndWorkflowID(ctx, organizationID, input.FromStepID, input.ToStepID, input.WorkflowID)
	if err != nil {
		return entity.Connexion{}, err
	}

	if len(connexions) > 0 {
		return entity.Connexion{}, nil
	}

	connexion := &entity.Connexion{
		ID:         uuid.New(),
		WorkflowID: input.WorkflowID,
		FromStepID: input.FromStepID,
		ToStepID:   input.ToStepID,
	}

	if err := (*u.connexionRepo).Create(ctx, connexion); err != nil {
		return entity.Connexion{}, err
	}

	return *connexion, nil
}
