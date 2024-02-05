package repository

import (
	"database/sql"
	"log"
	dbUtils "price-tracking-products/src/DB/utils"
	"price-tracking-products/src/models"
)

type ProductsRepository interface {
	AddProduct(models.Product) error
	DeleteProduct(models.Product) error
	AddProductToUser(models.User, models.Product) error
	RemoveProductToUser(models.User, models.Product) error
	UpdateProductPrice(models.Product) error
	ProductExists(string) (bool, error)
	GetProductByURL(string) (*models.Product, error)
}

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo() *ProductRepo {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil
	}
	return &ProductRepo{db: db}
}

func (p ProductRepo) AddProduct(product models.Product) error {
	statement := `insert into "products" ("id","name","brand",
	"higherprice","lowerprice","otherprice","discount","imageurl",
	"store","producturl") values ($1, $2, $3 ,$4 ,$5 ,$6 ,$7 ,$8 ,$9 ,$10)`
	_, err := p.db.Exec(statement, product.Id, product.Name, product.Brand, product.HigherPrice, product.LowePrice, product.OtherPaymentLowerPrice, product.Discount, product.ImageURL, product.Store, product.ProductURL)
	if err != nil {
		log.Printf("Error adding %s product %s \n", product.Name, err.Error())
		return err
	}
	return nil
}

func (p ProductRepo) DeleteProduct(product models.Product) error {
	statement := `DELETE * FROM products WHERE products.id=$1`
	_, err := p.db.Exec(statement, product.Id)
	if err != nil {
		log.Printf("Error deleting the product product %s. %s\n", product.Name, err.Error())
		return err
	}
	return nil
}

func (p ProductRepo) AddProductToUser(user models.User, product models.Product) error {
	statement := `insert into "products_users" ("userid","productid") values ($1, $2)`
	_, err := p.db.Exec(statement, user.Id, product.Id)
	if err != nil {
		log.Printf("Error adding the product %s to the user %s user %s\n", product.Name, user.UserName, err.Error())
		return err
	}
	return nil
}

func (p ProductRepo) RemoveProductToUser(user models.User, product models.Product) error {
	statement := `DELETE FROM "products_users" WHERE "userid"=$1 AND "productid"=$2;`
	_, err := p.db.Exec(statement, user.Id, product.Id)
	if err != nil {
		log.Printf("Error removing the product %s to the user %s user. %s \n", product.Name, user.UserName, err.Error())
		return err
	}
	return nil
}

func (p ProductRepo) UpdateProductPrice(newProduct models.Product) error {
	statement := `UPDATE "products" 
	SET "higherprice"=$1,"lowerprice"=$2,"otherprice"=$3,"discount"=$4 WHERE products.id = $5;`
	_, err := p.db.Exec(statement, newProduct.HigherPrice, newProduct.LowePrice, newProduct.OtherPaymentLowerPrice, newProduct.Discount, newProduct.Id)
	if err != nil {
		log.Printf("Error updating the product %s %s  \n", newProduct.Name, err.Error())
		return err
	}
	return nil
}

func (p ProductRepo) ProductExists(url string) (bool, error) {
	statement := `SELECT COUNT(*) AS EXISTS FROM products 
	WHERE producturl=$1`
	rows, err := p.db.Query(statement, url)
	if err != nil {
		log.Printf("Error searching if product exists %s.\n", err.Error())
		return false, err
	}
	defer rows.Close()

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

func (p ProductRepo) GetProductByURL(url string) (*models.Product, error) {
	statement := `SELECT id,name,brand,higherprice,lowerprice,otherprice,discount,imageurl,producturl,store 
	FROM products 
	WHERE producturl=$1`
	rows, err := p.db.Query(statement, url)
	if err != nil {
		log.Printf("Error searching if product exists %s.\n", err.Error())
		return nil, err
	}
	defer rows.Close()

	product := &models.Product{}
	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Brand, &product.HigherPrice, &product.LowePrice, &product.OtherPaymentLowerPrice, &product.Discount, &product.ImageURL, &product.Store, &product.ProductURL)
		if err != nil {
			log.Println("Error scanning the query result", err)
			return nil, err
		}
	}

	return product, nil
}
