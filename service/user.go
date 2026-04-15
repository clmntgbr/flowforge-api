package service

import (
	"forgeflow-api/domain"
	"forgeflow-api/dto"
	"forgeflow-api/repository"
	"time"

)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) FindByClerkID(clerkID string) *domain.User {
	return s.userRepository.FindByClerkID(clerkID)
}

func (s *UserService) CreateUser(id string, firstName string, lastName string, banned bool) (*domain.User, error) {
	user := &domain.User{
		ClerkID:   id,
		FirstName: firstName,
		LastName:  lastName,
		Banned:    banned,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(id string, firstName string, lastName string, banned bool) error {
	user := s.userRepository.FindByClerkID(id)

	if user == nil {
		return nil
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Banned = banned

	return s.userRepository.Update(user)
}

func (s *UserService) DeleteUser(id string) error {
	user := s.userRepository.FindByClerkID(id)

	if user == nil {
		return nil
	}

	return s.userRepository.Delete(user)
}

func (s *UserService) GetUser(user *domain.User) (*dto.UserOutput, error) {
	output := dto.NewUserOutput(*user)
	return &output, nil
}
