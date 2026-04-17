package repository

import (
	"context"
	"forgeflow-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConnexionRepository struct {
	db *gorm.DB
}

func NewConnexionRepository(db *gorm.DB) *ConnexionRepository {
	return &ConnexionRepository{db: db}
}

func (r *ConnexionRepository) Create(connexion *domain.Connexion) error {
	return r.db.Create(connexion).Error
}

func (r *ConnexionRepository) Update(connexion *domain.Connexion) error {
	return r.db.Save(connexion).Error
}

func (r *ConnexionRepository) Delete(connexion *domain.Connexion) error {
	return r.db.Delete(connexion).Error
}

func (r *ConnexionRepository) FindByFromTo(ctx context.Context, fromStepID uuid.UUID, toStepID uuid.UUID) ([]domain.Connexion, error) {
	var connexions []domain.Connexion
	err := r.db.WithContext(ctx).Where("from_step_id = ? AND to_step_id = ?", fromStepID, toStepID).Find(&connexions).Error
	if err != nil {
		return nil, err
	}
	return connexions, nil
}

func (r *ConnexionRepository) FindByID(ctx context.Context, connexionID uuid.UUID) (*domain.Connexion, error) {
	var connexion domain.Connexion
	err := r.db.WithContext(ctx).Where("id = ?", connexionID).First(&connexion).Error
	if err != nil {
		return nil, err
	}
	return &connexion, nil
}
