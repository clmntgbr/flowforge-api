package repository

import (
	"errors"
	"forgeflow-api/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(user *domain.User) error {
	return r.db.Delete(user).Error
}

func (r *UserRepository) FindByClerkID(clerkID string) (*domain.User, error) {
	var user domain.User

	err := r.db.
		Preload("ActiveProject").
		Where("clerk_id = ?", clerkID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) *domain.User {
	var user domain.User
	err := r.db.Preload("ActiveProject").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}
