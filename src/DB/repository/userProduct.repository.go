package repository

type UserProductRepository interface {
	AddUserProduct(Product, User) error
	DeleteUserProduct(Product, User) error
}
