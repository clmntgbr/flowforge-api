package usecase

import (
	"context"
	"forgeflow-api/domain"
)

type CreateUserUsecase struct {
	users           UserProvisioner
	createProjectUC *CreateProjectUsecase
}

func NewCreateUserUsecase(users UserProvisioner, createProjectUC *CreateProjectUsecase) *CreateUserUsecase {
	return &CreateUserUsecase{
		users:           users,
		createProjectUC: createProjectUC,
	}
}

func (s *CreateUserUsecase) CreateUser(ctx context.Context, id string, firstName string, lastName string, banned bool) (*domain.User, *domain.Project, error) {
	if s.users.FindByClerkID(id) != nil {
		return nil, nil, nil
	}

	user, err := s.users.CreateUser(id, firstName, lastName, banned)
	if err != nil {
		return nil, nil, err
	}

	project, err := s.createProjectUC.CreateProject(ctx, user, "Default Project")
	if err != nil {
		return nil, nil, err
	}

	return user, project, nil
}
