package usecase

import (
	"forgeflow-api/domain"
	"forgeflow-api/dto"
)

type GetUserUsecase struct {
	presenter UserPresenter
}

func NewGetUserUsecase(presenter UserPresenter) *GetUserUsecase {
	return &GetUserUsecase{presenter: presenter}
}

func (u *GetUserUsecase) GetUser(user *domain.User) (*dto.UserOutput, error) {
	return u.presenter.GetUser(user)
}
