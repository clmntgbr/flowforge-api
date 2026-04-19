package repository

import "forgeflow-api/domain"

type UserRepositoryInterface interface {
	FindByClerkID(clerkID string) (*domain.User, error)
	Update(user *domain.User) error
}
