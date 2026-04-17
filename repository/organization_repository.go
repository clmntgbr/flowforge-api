package repository

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(organization *domain.Organization) error {
	return r.db.Create(organization).Error
}

func (r *OrganizationRepository) Update(organization *domain.Organization) error {
	return r.db.Save(organization).Error
}

func (r *OrganizationRepository) Delete(organization *domain.Organization) error {
	return r.db.Delete(organization).Error
}

func (r *OrganizationRepository) CountOrganizationsByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Organization{}).
		Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("user_organizations.user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OrganizationRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Organization, error) {
	var organizations []domain.Organization

	db := r.db.WithContext(ctx).
		Model(&domain.Organization{}).
		Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("user_organizations.user_id = ?", userID)

	err := db.Find(&organizations).Error
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (r *OrganizationRepository) FindByUserIDAndOrganizationID(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (*domain.Organization, error) {
	var organization domain.Organization

	err := r.db.WithContext(ctx).
		Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("organizations.id = ? AND user_organizations.user_id = ?", organizationID, userID).
		First(&organization).Error

	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) ActivateOrganization(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (*domain.Organization, error) {
	var organization domain.Organization
	err := r.db.WithContext(ctx).
		Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("organizations.id = ? AND user_organizations.user_id = ?", organizationID, userID).
		First(&organization).Error

	if err != nil {
		return nil, errors.ErrOrganizationNotFound
	}

	err = r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Update("active_organization_id", organizationID).Error

	if err != nil {
		return nil, err
	}

	return &organization, nil
}
