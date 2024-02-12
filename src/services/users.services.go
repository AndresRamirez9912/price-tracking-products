package services

import (
	"price-tracking-products/src/DB/repository"
	"price-tracking-products/src/models"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService() (*UserService, error) {
	repo, err := repository.NewUserRepo()
	if err != nil {
		return nil, err
	}
	return &UserService{repo: repo}, nil
}

func (u UserService) AddUser(user models.User) error {
	err := u.repo.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) DeleteUser(user models.User) error {
	err := u.repo.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) ListUserProducts(user models.User) ([]models.Product, error) {
	products, err := u.repo.ListUserProducts(user)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (u UserService) HaveProduct(user models.User, url string) (bool, error) {
	haveProducts, err := u.repo.HaveProduct(user, url)
	if err != nil {
		return false, err
	}
	return haveProducts, nil
}
