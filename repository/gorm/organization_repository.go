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
