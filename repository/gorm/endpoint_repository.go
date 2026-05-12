package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"
	"flowforge-api/infrastructure/paginate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type endpointRepository struct {
	db *gorm.DB
}

func NewEndpointRepository(db *gorm.DB) repository.EndpointRepository {
	return &endpointRepository{db: db}
}

func (r *endpointRepository) Create(ctx context.Context, endpoint *entity.Endpoint) error {
	return r.db.WithContext(ctx).Create(endpoint).Error
}

func (r *endpointRepository) Update(ctx context.Context, endpoint *entity.Endpoint) error {
	return r.db.WithContext(ctx).Save(endpoint).Error
}

func (r *endpointRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Endpoint{}, id).Error
}

func (r *endpointRepository) List(ctx context.Context, organizationID uuid.UUID, query paginate.PaginateQuery) ([]entity.Endpoint, int64, error) {
	var endpoints []entity.Endpoint

	db := r.db.WithContext(ctx).Model(&entity.Endpoint{}).
		Where("organization_id = ?", organizationID)

	if query.Search != "" {
		db = db.Where("name ILIKE ?", "%"+query.Search+"%")
	}

	db, total, err := Paginate(db, query)
	if err != nil {
		return nil, 0, err
	}

	err = db.Find(&endpoints).Error
	if err != nil {
		return nil, 0, err
	}

	return endpoints, total, nil
}
