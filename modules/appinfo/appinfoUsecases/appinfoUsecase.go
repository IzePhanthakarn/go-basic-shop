package appinfoUsecases

import (
	"github.com/IzePhanthakarn/kawaii-shop/modules/appinfo"
	"github.com/IzePhanthakarn/kawaii-shop/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
}

type appinfoUsecase struct {
	appinfoRepository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.appinfoRepository.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	err := u.appinfoRepository.InsertCategory(req)
	if err != nil {
		return err
	}

	return nil
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	err := u.appinfoRepository.DeleteCategory(categoryId)
	if err != nil {
		return err
	}
	return nil
}