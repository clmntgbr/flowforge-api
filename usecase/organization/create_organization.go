package organization

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/presenter"
)

type CreateOrganizationUseCase struct {
	organizationRepo repository.OrganizationRepository
}

func NewCreateOrganizationUseCase(organizationRepo repository.OrganizationRepository) *CreateOrganizationUseCase {
	return &CreateOrganizationUseCase{organizationRepo: organizationRepo}
}

func (u *CreateOrganizationUseCase) Execute(ctx context.Context, user *entity.User, name string) (presenter.OrganizationResponse, error) {
	organization := &entity.Organization{
		Name: name,
		Users: []entity.User{
			{
				ID: user.ID,
			},
		},
	}

	if err := u.organizationRepo.Create(ctx, organization); err != nil {
		return presenter.OrganizationResponse{}, err
	}

	activeID := organization.ID
	if user.ActiveOrganizationID != nil {
		activeID = *user.ActiveOrganizationID
	}

	return presenter.NewOrganizationResponse(*organization, activeID), nil
}
