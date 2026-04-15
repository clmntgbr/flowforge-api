package usecase

import "context"

type UpdateUserUsecase struct {
	users UserUpdater
}

func NewUpdateUserUsecase(users UserUpdater) *UpdateUserUsecase {
	return &UpdateUserUsecase{users: users}
}

func (u *UpdateUserUsecase) UpdateUser(_ context.Context, id string, firstName string, lastName string, banned bool) error {
	return u.users.UpdateUser(id, firstName, lastName, banned)
}
