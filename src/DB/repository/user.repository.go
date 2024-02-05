package repository

import (
	"database/sql"
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
}

func NewUserRepo() *UserRepo {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil
	}
	return &UserRepo{db: db}
}

func (u UserRepo) AddUser(user models.User) error {
	statement := `insert into "users" ("id","email","username") values ($1, $2, $3)`
	_, err := u.db.Exec(statement, user.Id, user.Email, user.UserName)
	if err != nil {
		log.Printf("Error adding %s user %s\n", user.UserName, err.Error())
		return err
	}
	return nil
}

func (u UserRepo) DeleteUser(user models.User) error {
	statement := `DELETE * INTO "users" WHERE id=&1`
	_, err := u.db.Exec(statement, user.Id)
	if err != nil {
		log.Printf("Error deleting %s user \n", user.UserName)
		return err
	}
	return nil
}

func (u UserRepo) ListUserProducts(user models.User) ([]models.Product, error) {
	statement := `SELECT * from "products" 
	INNER JOIN "products_users" ON products_user.productsid = product.id 
	WHERE products_users.userid = $1`
	rows, err := u.db.Query(statement, user.Id)
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

	return products, nil
}

func (u UserRepo) HaveProduct(user models.User, url string) (bool, error) {
	statement := `SELECT count(userid) FROM products_users 
	INNER JOIN products ON products_users.productid = products.id
	WHERE userid=$1 AND producturl=$2;`
	rows, err := u.db.Query(statement, user.Id, url)
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
	return false, nil
}
