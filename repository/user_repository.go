package repository

import (
	"forgeflow-api/domain"

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

func (r *UserRepository) FindByClerkID(clerkID string) *domain.User {
	var user domain.User
	err := r.db.Preload("ActiveProject").Where("clerk_id = ?", clerkID).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}
