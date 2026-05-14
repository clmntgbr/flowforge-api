package gorm

import (
	"context"
	"errors"
	"flowforge-api/domain/entity"
	"flowforge-api/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByClerkID(ctx context.Context, clerkID string) (*entity.User, error) {
	var user entity.User
	err := dbWithContext(ctx, r.db).Where("clerk_id = ?", clerkID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return dbWithContext(ctx, r.db).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return dbWithContext(ctx, r.db).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return dbWithContext(ctx, r.db).Delete(&entity.User{}, id).Error
}

func (r *userRepository) DeleteByClerkID(ctx context.Context, clerkID string) error {
	return dbWithContext(ctx, r.db).Where("clerk_id = ?", clerkID).Delete(&entity.User{}).Error
}
