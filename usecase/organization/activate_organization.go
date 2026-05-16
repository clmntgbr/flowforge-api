package organization

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type ActivateOrganizationUseCase struct {
	organizationRepo *repository.OrganizationRepository
}

func NewActivateOrganizationUseCase(organizationRepo *repository.OrganizationRepository) *ActivateOrganizationUseCase {
	return &ActivateOrganizationUseCase{organizationRepo: organizationRepo}
}

func (u *ActivateOrganizationUseCase) Execute(ctx context.Context, user *entity.User, organizationID uuid.UUID) (entity.Organization, error) {
	organization, err := (*u.organizationRepo).ActivateOrganization(ctx, user.ID, organizationID)
	if err != nil {
		return entity.Organization{}, err
	}

	return organization, nil
}
