package repository

import (
	"log"
	dbUtils "price-tracking-products/src/DB/utils"
)

type Product struct {
	Name                   string `json:"name"`
	Id                     string `json:"id"`
	Brand                  string `json:"brand"`
	HigherPrice            string `json:"higherPrice"`
	LowePrice              string `json:"lowePrice"`
	OtherPaymentLowerPrice string `json:"otherPaymentLowerPrice"`
	Discount               string `json:"discount"`
	ImageURL               string `json:"imageURL"`
	Store                  string `json:"store"`
	ProductURL             string `json:"productURL"`
}

type ProductsRepository interface {
	AddProduct() error
	DeleteProduct() error
	AddProductToUser(product Product) error
	RemoveProductToUser(product Product) error
}

func (p Product) AddProduct() error {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `insert into "products" ("id","name","brand",
	"higherprice","lowerprice","otherprice","discount","imageurl",
	"store","producturl") values ($1, $2, $3 ,$4 ,$5 ,$6 ,$7 ,$8 ,$9 ,$10)`
	_, err = db.Exec(statement, p.Id, p.Name, p.Brand, p.HigherPrice, p.LowePrice, p.OtherPaymentLowerPrice, p.Discount, p.ImageURL, p.Store, p.ProductURL)
	if err != nil {
		log.Printf("Error adding %s product %s \n", p.Name, err.Error())
		return err
	}
	return nil
}

func (p Product) AddProductToUser(user User) error {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `insert into "products_users" ("userid","productid") values ($1, $2)`
	_, err = db.Exec(statement, user.Id, p.Id)
	if err != nil {
		log.Printf("Error adding the product %s to the user %s user %s\n", p.Name, user.UserName, err.Error())
		return err
	}
	return nil
}

func (p Product) RemoveProductToUser(user User) error {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return err
	}
	defer dbUtils.CloseDBConnection(db)

	statement := `DELETE FROM "products_users" WHERE "userid"=$1 AND "productid"=$2;`
	_, err = db.Exec(statement, user.Id, p.Id)
	if err != nil {
		log.Printf("Error removing the product %s to the user %s user \n", p.Name, user.UserName)
		return err
	}
	return nil
}
