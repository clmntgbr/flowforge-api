package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) repository.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, organization *entity.Organization) error {
	return r.db.WithContext(ctx).Create(organization).Error
}

func (r *organizationRepository) Update(ctx context.Context, organization *entity.Organization) error {
	return r.db.WithContext(ctx).Save(organization).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Organization{}, id).Error
}

func (r *organizationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Organization, error) {
	var organizations []entity.Organization

	db := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("user_organizations.user_id = ?", userID)

	err := db.Find(&organizations).Error
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (r *organizationRepository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (entity.Organization, error) {
	var organization entity.Organization

	db := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("organizations.id = ? AND user_organizations.user_id = ?", id, userID)

	err := db.First(&organization).Error
	if err != nil {
		return entity.Organization{}, err
	}

	return organization, nil
}
