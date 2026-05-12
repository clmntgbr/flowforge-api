package organization

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
)

type ListOrganizationsUseCase struct {
	organizationRepo repository.OrganizationRepository
}

func NewListOrganizationsUseCase(organizationRepo repository.OrganizationRepository) *ListOrganizationsUseCase {
	return &ListOrganizationsUseCase{organizationRepo: organizationRepo}
}

func (u *ListOrganizationsUseCase) Execute(ctx context.Context, user *entity.User, activeOrganizationID uuid.UUID) ([]entity.Organization, error) {
	organizations, err := u.organizationRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return []entity.Organization{}, err
	}

	for _, organization := range organizations {
		if organization.ID == activeOrganizationID {
			organization.IsActive = true
		}
	}

	return organizations, nil
}
