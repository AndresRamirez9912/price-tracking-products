package repository

import (
	"log"
	dbUtils "price-tracking-products/src/DB/utils"
)

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	UserName string `json:"userName"`
}

type UserRepository interface {
	AddUser() error
	DeleteUser() error
	ListUserProducts() ([]Product, error)
	HaveProduct(Product) (bool, error)
}

func (u User) AddUser() error {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `insert into "users" ("id","email","username") values ($1, $2, $3)`
	_, err = db.Exec(statement, u.Id, u.Email, u.UserName)
	if err != nil {
		log.Printf("Error adding %s user %s\n", u.UserName, err.Error())
		return err
	}
	return nil
}

func (u User) DeleteUser() error {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `DELETE * INTO "users" WHERE id=&1`
	_, err = db.Exec(statement, u.Id)
	if err != nil {
		log.Printf("Error deleting %s user \n", u.UserName)
		return err
	}
	return nil
}

func (u User) ListUserProducts() ([]Product, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil, err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `SELECT * from "products" 
	INNER JOIN "products_users" ON products_user.productsid = product.id 
	WHERE products_users.userid = $1`
	rows, err := db.Query(statement, u.Id)
	if err != nil {
		log.Printf("Error listing the products of the user %s user. %s \n", u.UserName, err.Error())
		return nil, err
	}

	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
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

func (u User) HaveProduct(url string) (bool, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return false, err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `SELECT count(userid) FROM products_users 
	INNER JOIN products ON products_users.productid = products.id
	WHERE userid=$1 AND producturl=$2;`
	rows, err := db.Query(statement, u.Id, url)
	if err != nil {
		log.Printf("Error checking if user %s has the product. %s \n", u.UserName, err.Error())
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

// TODO: Move this method to the products or refactor
func (u User) ProductExists(url string) (bool, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return false, err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `SELECT COUNT(*) AS EXISTS FROM products 
	WHERE producturl=$1`
	rows, err := db.Query(statement, url)
	if err != nil {
		log.Printf("Error searching if product exists %s.\n", err.Error())
		return false, err
	}
	defer rows.Close()

	// TODO: Create a util to check the result and return a bool result
	var amount int
	for rows.Next() {
		err = rows.Scan(&amount)
		if err != nil {
			log.Println("Error scanning the query result", err)
			return false, err
		}
	}

	if amount != 0 {
		return true, nil
	}

	return false, nil
}

// TODO: Move this method to the products or refactor
func (u User) GetProductByURL(url string) (*Product, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil, err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `SELECT id,name,brand,higherprice,lowerprice,otherprice,discount,imageurl,producturl,store 
	FROM products 
	WHERE producturl=$1`
	rows, err := db.Query(statement, url)
	if err != nil {
		log.Printf("Error searching if product exists %s.\n", err.Error())
		return nil, err
	}
	defer rows.Close()

	product := &Product{}
	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Brand, &product.HigherPrice, &product.LowePrice, &product.OtherPaymentLowerPrice, &product.Discount, &product.ImageURL, &product.Store, &product.ProductURL)
		if err != nil {
			log.Println("Error scanning the query result", err)
			return nil, err
		}
	}

	return product, nil
}
