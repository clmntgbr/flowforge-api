package organization

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type GetOrganizationByIDUseCase struct {
	organizationRepo repository.OrganizationRepository
}

func NewGetOrganizationByIDUseCase(organizationRepo repository.OrganizationRepository) *GetOrganizationByIDUseCase {
	return &GetOrganizationByIDUseCase{organizationRepo: organizationRepo}
}

func (u *GetOrganizationByIDUseCase) Execute(ctx context.Context, user *entity.User, organizationID uuid.UUID) (entity.Organization, error) {
	organization, err := u.organizationRepo.GetByIDAndUserID(ctx, organizationID, user.ID)
	if err != nil {
		return entity.Organization{}, err
	}

	return organization, nil
}
