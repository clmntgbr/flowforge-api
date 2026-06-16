package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	return dbWithContext(ctx, r.db).Create(tag).Error
}

func (r *tagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	return dbWithContext(ctx, r.db).Save(tag).Error
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.Tag{}, id).Error
}

func (r *tagRepository) List(ctx context.Context, organizationID uuid.UUID) ([]entity.Tag, error) {
	var tags []entity.Tag

	err := dbWithContext(ctx, r.db).Model(&entity.Tag{}).
		Where("organization_id = ?", organizationID).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}
