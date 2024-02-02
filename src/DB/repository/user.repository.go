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
	AddUser(User) error
	DeleteUser(User) error
	ListUserProducts(User) ([]Product, error)
}

func (u User) AddUser(user User) error {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `insert into "users" ("id","email","userName") values ($1, $2, $3)`
	_, err = db.Exec(statement, user.Id, user.Email, user.UserName)
	if err != nil {
		log.Printf("Error adding %s user \n", user.UserName)
		return err
	}
	return nil
}

func (u User) DeleteUser(user User) error {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `DELETE * FROM "users" WHERE id=&1`
	_, err = db.Exec(statement, user.Id)
	if err != nil {
		log.Printf("Error deleting %s user \n", user.UserName)
		return err
	}
	return nil
}

func (u User) ListUserProducts(User) ([]Product, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil, err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `SELECT * from "products" 
	INNER JOIN "product_user" ON product_user.productsId = products.id 
	WHERE product_user.userId = $1`
	rows, err := db.Query(statement, u.Id)
	if err != nil {
		log.Printf("Error deleting %s user \n", u.UserName)
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
