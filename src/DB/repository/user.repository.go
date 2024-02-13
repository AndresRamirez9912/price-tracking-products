package repository

import (
	"database/sql"
	"errors"
	"log"
	dbUtils "price-tracking-products/src/DB/utils"
	"price-tracking-products/src/models"
)

type UserRepository interface {
	AddUser(models.User) error
	DeleteUser(models.User) error
	ListUserProducts(models.User) ([]models.Product, error)
	HaveProduct(models.User, string) (bool, error)
}

type UserRepo struct {
	db *sql.DB
	tx *sql.Tx
}

func NewUserRepo() (*UserRepo, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil, err
	}
	tx, err := dbUtils.CreateTransaction(db)
	if err != nil {
		return nil, err
	}
	return &UserRepo{db: db, tx: tx}, nil
}

func (u *UserRepo) updateTransaction() error {
	tx, err := dbUtils.CreateTransaction(u.db)
	if err != nil {
		return nil
	}
	u.tx = tx
	return nil
}

func (u UserRepo) AddUser(user models.User) error {
	// Update the transaction
	err := u.updateTransaction()
	if err != nil {
		return err
	}

	statement := `insert into "users" ("id","email","user_name") values ($1, $2, $3)`
	_, err = u.tx.Exec(statement, user.Id, user.Email, user.UserName)
	if err != nil {
		log.Printf("Error adding %s user %s\n", user.UserName, err.Error())
		return err
	}
	defer dbUtils.CloseTransaction(u.tx, err)
	return nil
}

func (u UserRepo) DeleteUser(user models.User) error {
	// Update the transaction
	err := u.updateTransaction()
	if err != nil {
		return err
	}

	statement := `DELETE * INTO "users" WHERE id=&1`
	_, err = u.tx.Exec(statement, user.Id)
	if err != nil {
		log.Printf("Error deleting %s user \n", user.UserName)
		return err
	}
	defer dbUtils.CloseTransaction(u.tx, err)
	return nil
}

func (u UserRepo) ListUserProducts(user models.User) ([]models.Product, error) {
	// Update the transaction
	err := u.updateTransaction()
	if err != nil {
		return nil, err
	}

	if (user == models.User{}) {
		log.Println("User empty, bad request")
		return nil, errors.New("User empty, bad request")
	}

	statement := `SELECT id,name,brand,higher_price,lower_price,other_price,discount,image_url,product_url,store from "products" 
	INNER JOIN "products_users" ON products_users.product_id = products.id 
	WHERE products_users.user_id = $1`
	rows, err := u.tx.Query(statement, user.Id)
	if err != nil {
		log.Printf("Error listing the products of the user %s user. %s \n", user.UserName, err.Error())
		return nil, err
	}

	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err = rows.Scan(&p.Id, &p.Name, &p.Brand, &p.HigherPrice, &p.LowePrice, &p.OtherPaymentLowerPrice, &p.Discount, &p.ImageURL, &p.ProductURL, &p.Store)
		if err != nil {
			log.Println("Error scanning the elements for the user", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error during the  iteration", err)
		return nil, err
	}
	defer dbUtils.CloseTransaction(u.tx, err)
	return products, nil
}

func (u UserRepo) HaveProduct(user models.User, url string) (bool, error) {
	// Update the transaction
	err := u.updateTransaction()
	if err != nil {
		return false, err
	}

	statement := `SELECT count(user_id) FROM products_users 
	INNER JOIN products ON products_users.product_id = products.id
	WHERE user_id=$1 AND product_url=$2;`
	rows, err := u.tx.Query(statement, user.Id, url)
	if err != nil {
		log.Printf("Error checking if user %s has the product. %s \n", user.UserName, err.Error())
		return false, err
	}
	defer rows.Close()

	var amount int
	for rows.Next() {
		err = rows.Scan(&amount)
		if err != nil {
			log.Println("Error scanning the element for the user", err)
			return false, err
		}
	}

	if amount != 0 {
		return true, nil
	}
	defer dbUtils.CloseTransaction(u.tx, err)
	return false, nil
}
