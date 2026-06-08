package connexion

import (
	"context"
	"flowforge-api/domain/repository"
	usecaseStep "flowforge-api/usecase/step"

	"github.com/google/uuid"
)

type DeleteConnexionUseCase struct {
	connexionRepo     *repository.ConnexionRepository
	assignTreeIndices *usecaseStep.AssignTreeIndicesUseCase
}

func NewDeleteConnexionUseCase(connexionRepo *repository.ConnexionRepository, assignTreeIndices *usecaseStep.AssignTreeIndicesUseCase) *DeleteConnexionUseCase {
	return &DeleteConnexionUseCase{connexionRepo: connexionRepo, assignTreeIndices: assignTreeIndices}
}

func (u *DeleteConnexionUseCase) Execute(ctx context.Context, organizationID uuid.UUID, id uuid.UUID) error {
	conn, err := (*u.connexionRepo).GetByIDAndOrganizationID(ctx, organizationID, id)
	if err != nil {
		return err
	}

	if err := (*u.connexionRepo).Delete(ctx, conn.ID); err != nil {
		return err
	}

	return u.assignTreeIndices.Execute(ctx, conn.WorkflowID)
}
