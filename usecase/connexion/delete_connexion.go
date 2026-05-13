package connexion

import (
	"context"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type DeleteConnexionUseCase struct {
	connexionRepo repository.ConnexionRepository
}

func NewDeleteConnexionUseCase(connexionRepo repository.ConnexionRepository) *DeleteConnexionUseCase {
	return &DeleteConnexionUseCase{connexionRepo: connexionRepo}
}

func (u *DeleteConnexionUseCase) Execute(ctx context.Context, organizationID uuid.UUID, id uuid.UUID) error {
	connexion, err := u.connexionRepo.GetByIDAndOrganizationID(ctx, organizationID, id)
	if err != nil {
		return err
	}

	if err := u.connexionRepo.Delete(ctx, connexion.ID); err != nil {
		return err
	}

	return nil
}
