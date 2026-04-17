package service

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/repository"
	"forgeflow-api/rules"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OrganizationService struct {
	organizationRepository *repository.OrganizationRepository
	organizationRules      *rules.OrganizationRules
}

func NewOrganizationService(organizationRepository *repository.OrganizationRepository, organizationRules *rules.OrganizationRules) *OrganizationService {
	return &OrganizationService{
		organizationRepository: organizationRepository,
		organizationRules:      organizationRules,
	}
}

func (s *OrganizationService) CreateOrganization(c fiber.Ctx, user *domain.User, name string) (dto.OrganizationOutput, error) {
	if err := s.organizationRules.MaxOrganizationsPerUser(c.Context(), user.ID); err != nil {
		return dto.OrganizationOutput{}, err
	}

	organization := &domain.Organization{
		Name: name,
		Users: []domain.User{
			{
				ID: user.ID,
			},
		},
	}

	if err := s.organizationRepository.Create(organization); err != nil {
		return dto.OrganizationOutput{}, err
	}

	activeID := uuid.Nil
	activeID = organization.ID

	if user.ActiveOrganizationID != nil {
		activeID = *user.ActiveOrganizationID
	}

	return dto.NewOrganizationOutput(*organization, activeID), nil
}

func (s *OrganizationService) GetOrganizations(c fiber.Ctx, user *domain.User, activeOrganizationID uuid.UUID) ([]dto.MinimalOrganizationOutput, error) {
	organizations, err := s.organizationRepository.FindAllByUserID(c.Context(), user.ID)
	if err != nil {
		return nil, err
	}

	return dto.NewMinimalOrganizationsOutput(organizations, activeOrganizationID), nil
}

func (s *OrganizationService) GetOrganizationByID(c fiber.Ctx, user *domain.User, organizationUUID uuid.UUID) (dto.OrganizationOutput, error) {
	organization, err := s.organizationRepository.FindByUserIDAndOrganizationID(c.Context(), organizationUUID, user.ID)
	if err != nil {
		return dto.OrganizationOutput{}, err
	}

	return dto.NewOrganizationOutput(*organization, organization.ID), nil
}

func (s *OrganizationService) UpdateOrganization(c fiber.Ctx, user *domain.User, organizationUUID uuid.UUID, req dto.UpdateOrganizationInput) (dto.OrganizationOutput, error) {
	organization, err := s.organizationRepository.FindByUserIDAndOrganizationID(c.Context(), organizationUUID, user.ID)
	if err != nil {
		return dto.OrganizationOutput{}, err
	}

	organization.Name = req.Name
	if err := s.organizationRepository.Update(organization); err != nil {
		return dto.OrganizationOutput{}, err
	}

	return dto.NewOrganizationOutput(*organization, organization.ID), nil
}

func (s *OrganizationService) ActivateOrganization(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (dto.OrganizationOutput, error) {
	organization, err := s.organizationRepository.ActivateOrganization(ctx, userID, organizationID)
	if err != nil {
		return dto.OrganizationOutput{}, err
	}

	return dto.NewOrganizationOutput(*organization, organization.ID), nil
}
