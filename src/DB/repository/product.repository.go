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
	db *sql.Tx
}

func NewProductRepo() *ProductRepo {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil
	}
	tx, err := dbUtils.CreateTransaction(db)
	if err != nil {
		return nil
	}
	return &ProductRepo{db: tx}
}

func (p ProductRepo) AddProduct(product models.Product) error {
	statement := `insert into "products" ("id","name","brand",
	"higher_price","lower_price","other_price","discount","image_url",
	"store","product_url") values ($1, $2, $3 ,$4 ,$5 ,$6 ,$7 ,$8 ,$9 ,$10)`
	_, err := p.db.Exec(statement, product.Id, product.Name, product.Brand, product.HigherPrice, product.LowePrice, product.OtherPaymentLowerPrice, product.Discount, product.ImageURL, product.Store, product.ProductURL)
	if err != nil {
		log.Printf("Error adding %s product %s \n", product.Name, err.Error())
		return err
	}
	defer dbUtils.CloseTransaction(p.db, err)
	return nil
}

func (p ProductRepo) DeleteProduct(product models.Product) error {
	statement := `DELETE * FROM products WHERE products.id=$1`
	_, err := p.db.Exec(statement, product.Id)
	if err != nil {
		log.Printf("Error deleting the product product %s. %s\n", product.Name, err.Error())
		return err
	}
	defer dbUtils.CloseTransaction(p.db, err)
	return nil
}

func (p ProductRepo) AddProductToUser(user models.User, product models.Product) error {
	statement := `insert into "products_users" ("user_id","product_id") values ($1, $2)`
	_, err := p.db.Exec(statement, user.Id, product.Id)
	if err != nil {
		log.Printf("Error adding the product %s to the user %s user %s\n", product.Name, user.UserName, err.Error())
		return err
	}
	defer dbUtils.CloseTransaction(p.db, err)
	return nil
}

func (p ProductRepo) RemoveProductToUser(user models.User, product models.Product) error {
	statement := `DELETE FROM "products_users" WHERE "user_id"=$1 AND "product_id"=$2;`
	_, err := p.db.Exec(statement, user.Id, product.Id)
	if err != nil {
		log.Printf("Error removing the product %s to the user %s user. %s \n", product.Name, user.UserName, err.Error())
		return err
	}
	defer dbUtils.CloseTransaction(p.db, err)
	return nil
}

func (p ProductRepo) UpdateProductPrice(newProduct models.Product) error {
	statement := `UPDATE "products" 
	SET "higher_price"=$1,"lower_price"=$2,"other_price"=$3,"discount"=$4 WHERE products.id = $5;`
	_, err := p.db.Exec(statement, newProduct.HigherPrice, newProduct.LowePrice, newProduct.OtherPaymentLowerPrice, newProduct.Discount, newProduct.Id)
	if err != nil {
		log.Printf("Error updating the product %s %s  \n", newProduct.Name, err.Error())
		return err
	}
	defer dbUtils.CloseTransaction(p.db, err)
	return nil
}

func (p ProductRepo) ProductExists(url string) (bool, error) {
	statement := `SELECT COUNT(*) AS EXISTS FROM products 
	WHERE product_url=$1`
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
	defer dbUtils.CloseTransaction(p.db, err)
	return false, nil
}

func (p ProductRepo) GetProductByURL(url string) (*models.Product, error) {
	statement := `SELECT id,name,brand,higher_price,lower_price,other_price,discount,image_url,product_url,store 
	FROM products 
	WHERE product_url=$1`
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
	defer dbUtils.CloseTransaction(p.db, err)
	return product, nil
}
