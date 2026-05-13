package gorm

import (
	"context"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type connexionRepository struct {
	db *gorm.DB
}

func NewConnexionRepository(db *gorm.DB) repository.ConnexionRepository {
	return &connexionRepository{db: db}
}

func (r *connexionRepository) Create(ctx context.Context, connexion *entity.Connexion) error {
	return r.db.WithContext(ctx).Create(connexion).Error
}

func (r *connexionRepository) Update(ctx context.Context, connexion *entity.Connexion) error {
	return r.db.WithContext(ctx).Save(connexion).Error
}

func (r *connexionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Connexion{}, id).Error
}

func (r *connexionRepository) GetByFromStepIDAndToStepIDAndWorkflowID(ctx context.Context, organizationID uuid.UUID, fromStepID uuid.UUID, toStepID uuid.UUID, workflowID uuid.UUID) ([]entity.Connexion, error) {
	var connexions []entity.Connexion

	err := r.db.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = connexions.workflow_id").
		Where("connexions.from_step_id = ? AND connexions.to_step_id = ? AND connexions.workflow_id = ? AND workflows.organization_id = ?", fromStepID, toStepID, workflowID, organizationID).
		Find(&connexions).Error

	if err != nil {
		return nil, err
	}
	return connexions, nil
}

func (r *connexionRepository) GetByIDAndOrganizationID(ctx context.Context, organizationID uuid.UUID, id uuid.UUID) (entity.Connexion, error) {
	var connexion entity.Connexion

	err := r.db.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = connexions.workflow_id").
		Where("connexions.id = ? AND workflows.organization_id = ?", id, organizationID).
		First(&connexion).Error
	if err != nil {
		return entity.Connexion{}, err
	}
	return connexion, nil
}
