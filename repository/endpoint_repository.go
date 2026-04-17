package repository

import (
	"context"
	"forgeflow-api/domain"
	"forgeflow-api/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EndpointRepository struct {
	db *gorm.DB
}

func NewEndpointRepository(db *gorm.DB) *EndpointRepository {
	return &EndpointRepository{db: db}
}

func (r *EndpointRepository) Create(endpoint *domain.Endpoint) error {
	return r.db.Create(endpoint).Error
}

func (r *EndpointRepository) Update(endpoint *domain.Endpoint) error {
	return r.db.Save(endpoint).Error
}

func (r *EndpointRepository) Delete(endpoint *domain.Endpoint) error {
	return r.db.Delete(endpoint).Error
}

func (r *EndpointRepository) FindAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, q dto.PaginateQuery) ([]domain.Endpoint, int64, error) {
	var endpoints []domain.Endpoint

	db := r.db.WithContext(ctx).Model(&domain.Endpoint{}).
		Where("organization_id = ?", organizationID)

	if q.Search != "" {
		db = db.Where("name ILIKE ?", "%"+q.Search+"%")
	}

	db, total, err := Paginate(db, q)
	if err != nil {
		return nil, 0, err
	}

	err = db.Find(&endpoints).Error
	if err != nil {
		return nil, 0, err
	}

	return endpoints, total, nil
}

func (r *EndpointRepository) FindByOrganizationIDAndEndpointID(ctx context.Context, organizationID uuid.UUID, endpointID uuid.UUID) (domain.Endpoint, error) {
	var endpoint domain.Endpoint
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND id = ?", organizationID, endpointID).
		First(&endpoint).Error
	if err != nil {
		return domain.Endpoint{}, err
	}
	return endpoint, nil
}

func (r *EndpointRepository) FindByID(ctx context.Context, endpointID uuid.UUID) (domain.Endpoint, error) {
	var endpoint domain.Endpoint
	err := r.db.WithContext(ctx).Where("id = ?", endpointID).First(&endpoint).Error
	if err != nil {
		return domain.Endpoint{}, err
	}
	return endpoint, nil
}
