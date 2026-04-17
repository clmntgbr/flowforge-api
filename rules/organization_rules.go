package rules

import (
	"context"
	"forgeflow-api/errors"
	"forgeflow-api/repository"

	"github.com/google/uuid"
)

type OrganizationRules struct {
	organizationRepo *repository.OrganizationRepository
}

func NewOrganizationRules(organizationRepo *repository.OrganizationRepository) *OrganizationRules {
	return &OrganizationRules{
		organizationRepo: organizationRepo,
	}
}
func (p *OrganizationRules) MaxOrganizationsPerUser(ctx context.Context, userID uuid.UUID) error {
	count, err := p.organizationRepo.CountOrganizationsByUserID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	if count >= 10 {
		return errors.ErrMaxOrganizationsReached
	}

	return nil
}
