package middlewaresUsecases

import (
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares"
	"github.com/IzePhanthakarn/kawaii-shop/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewareRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewareRepository: middlewareRepository,
	}
}

func (u *middlewaresUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middlewareRepository.FindAccessToken(userId, accessToken)
}

func (u *middlewaresUsecase) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.middlewareRepository.FindRole()
	if err != nil {
		return nil, err
	}

	return roles, nil
}
