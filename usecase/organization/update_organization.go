package organization

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type UpdateOrganizationUseCase struct {
	organizationRepo repository.OrganizationRepository
}

func NewUpdateOrganizationUseCase(organizationRepo repository.OrganizationRepository) *UpdateOrganizationUseCase {
	return &UpdateOrganizationUseCase{organizationRepo: organizationRepo}
}

func (u *UpdateOrganizationUseCase) Execute(ctx context.Context, user *entity.User, organizationID uuid.UUID, name string) (entity.Organization, error) {
	organization, err := u.organizationRepo.GetByIDAndUserID(ctx, organizationID, user.ID)
	if err != nil {
		return entity.Organization{}, err
	}

	organization.Name = name

	err = u.organizationRepo.Update(ctx, &organization)
	if err != nil {
		return entity.Organization{}, err
	}

	return organization, nil
}
