package services

import (
	"price-tracking-products/src/models"
	"price-tracking-products/src/models/repository"
)

type UserServiceInterface interface {
	AddUser(user models.User) error
	DeleteUser(user models.User) error
	ListUserProducts(user *models.User) ([]models.Product, error)
}

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{repo: userRepo}
}

func (u UserService) AddUser(user models.User) error {
	authResponse, err := u.repo.CreateUser(user)
	if err != nil {
		return err
	}

	user.Id = authResponse.Response.UserSub
	err = u.repo.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) DeleteUser(user models.User) error {
	// Delete the products linked
	err := u.repo.DeleteAllUserProducts(user)
	if err != nil {
		return err
	}

	// Delete the user
	err = u.repo.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) ListUserProducts(user *models.User) ([]models.Product, error) {
	products, err := u.repo.ListUserProducts(*user)
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
