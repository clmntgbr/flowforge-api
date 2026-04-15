package usecase

import "context"

type DeleteUserUsecase struct {
	users UserDeleter
}

func NewDeleteUserUsecase(users UserDeleter) *DeleteUserUsecase {
	return &DeleteUserUsecase{users: users}
}

func (u *DeleteUserUsecase) DeleteUser(_ context.Context, id string) error {
	return u.users.DeleteUser(id)
}
