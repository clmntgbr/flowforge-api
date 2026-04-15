package usecase

import (
	"context"
	"errors"
	"forgeflow-api/domain"
	apperrors "forgeflow-api/errors"
	"forgeflow-api/repository"
)

type CreateUserUsecase struct {
	users                UserProvisioner
	createProjectUsecase *CreateProjectUsecase
	userRepo             *repository.UserRepository
}

func NewCreateUserUsecase(users UserProvisioner, createProjectUsecase *CreateProjectUsecase, userRepo *repository.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{
		users:                users,
		createProjectUsecase: createProjectUsecase,
		userRepo:             userRepo,
	}
}

func (s *CreateUserUsecase) CreateUser(ctx context.Context, id string, firstName string, lastName string, banned bool) (*domain.User, *domain.Project, error) {
	existing, err := s.users.FindByClerkID(id)
	if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
		return nil, nil, err
	}
	if existing != nil {
		return nil, nil, nil
	}

	user, err := s.users.CreateUser(id, firstName, lastName, banned)
	if err != nil {
		return nil, nil, err
	}

	project, err := s.createProjectUsecase.CreateProject(ctx, user, "Default Project")
	if err != nil {
		return nil, nil, err
	}

	user.ActiveProject = project
	user.ActiveProjectID = &project.ID

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, nil, err
	}

	return user, project, nil
}
