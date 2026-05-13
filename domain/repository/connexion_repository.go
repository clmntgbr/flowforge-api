package repository

import (
	"context"
	"flowforge-api/domain/entity"

	"github.com/google/uuid"
)

type ConnexionRepository interface {
	Create(ctx context.Context, connexion *entity.Connexion) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByFromStepIDAndToStepIDAndWorkflowID(ctx context.Context, organizationID uuid.UUID, fromStepID uuid.UUID, toStepID uuid.UUID, workflowID uuid.UUID) ([]entity.Connexion, error)
	GetByIDAndOrganizationID(ctx context.Context, organizationID uuid.UUID, id uuid.UUID) (entity.Connexion, error)
}
