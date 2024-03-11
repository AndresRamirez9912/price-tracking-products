package repository

import (
	"database/sql"
	"log"
	"price-tracking-products/src/models"
	dbUtils "price-tracking-products/src/models/utils"
)

type ProductsRepository interface {
	AddProduct(models.Product) error
	DeleteProduct(models.Product) error
	AddProductToUser(models.User, models.Product) error
	RemoveProductToUser(models.User, models.Product) error
	UpdateProductPrice(models.Product) error
	ProductExists(string) (bool, error)
	GetProductByURL(string) (*models.Product, error)

	AddProductHistory(*models.Product) error
	DeleteProductHistory(*models.Product) error
	GetProductHistory(*models.Product) ([]models.ProductHistory, error)
}

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo() (*ProductRepo, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil, err
	}
	return &ProductRepo{db: db}, nil
}

func (repo ProductRepo) AddProduct(product models.Product) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `insert into "products" ("id","name","brand",
	"higher_price","lower_price","other_price","discount","image_url",
	"store","product_url") values ($1, $2, $3 ,$4 ,$5 ,$6 ,$7 ,$8 ,$9 ,$10)`
	_, err = tx.Exec(statement, product.Id, product.Name, product.Brand, product.HigherPrice, product.LowePrice, product.OtherPaymentLowerPrice, product.Discount, product.ImageURL, product.Store, product.ProductURL)
	if err != nil {
		log.Printf("Error adding %s product %s \n", product.Name, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo ProductRepo) DeleteProduct(product models.Product) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `DELETE * FROM products WHERE products.id=$1`
	_, err = tx.Exec(statement, product.Id)
	if err != nil {
		log.Printf("Error deleting the product product %s. %s\n", product.Name, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo ProductRepo) AddProductToUser(user models.User, product models.Product) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `insert into "products_users" ("user_id","product_id") values ($1, $2)`
	_, err = tx.Exec(statement, user.Id, product.Id)
	if err != nil {
		log.Printf("Error adding the product %s to the user %s. %s\n", product.Name, user.UserName, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo ProductRepo) RemoveProductToUser(user models.User, product models.Product) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `DELETE FROM "products_users" WHERE "user_id"=$1 AND "product_id"=$2;`
	_, err = tx.Exec(statement, user.Id, product.Id)
	if err != nil {
		log.Printf("Error removing the product %s to the user %s user. %s \n", product.Name, user.UserName, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo ProductRepo) UpdateProductPrice(newProduct models.Product) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `UPDATE "products" 
	SET "higher_price"=$1,"lower_price"=$2,"other_price"=$3,"discount"=$4 WHERE products.id = $5;`
	_, err = tx.Exec(statement, newProduct.HigherPrice, newProduct.LowePrice, newProduct.OtherPaymentLowerPrice, newProduct.Discount, newProduct.Id)
	if err != nil {
		log.Printf("Error updating the product %s %s  \n", newProduct.Name, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo ProductRepo) ProductExists(url string) (bool, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `SELECT COUNT(*) AS EXISTS FROM products 
	WHERE product_url=$1`
	rows, err := tx.Query(statement, url)
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

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return false, nil
}

func (repo ProductRepo) GetProductByURL(url string) (*models.Product, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `SELECT id,name,brand,higher_price,lower_price,other_price,discount,image_url,product_url,store 
	FROM products 
	WHERE product_url=$1`
	rows, err := tx.Query(statement, url)
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

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repo ProductRepo) AddProductHistory(product *models.Product) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `insert into "product_history" 
	("product_id","product_name","higher_price","lower_price","other_price") 
	values ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(statement, product.Id, product.Name, product.HigherPrice, product.LowePrice, product.OtherPaymentLowerPrice)
	if err != nil {
		log.Printf("Error adding the product %s to the history table %s\n", product.Name, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (repo ProductRepo) DeleteProductHistory(product *models.Product) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `DELETE * FROM product_history WHERE product_history.id=$1`
	_, err = tx.Exec(statement, product.Id)
	if err != nil {
		log.Printf("Error deleting the product %s to the history table %s\n", product.Name, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo ProductRepo) GetProductHistory(product *models.Product) ([]models.ProductHistory, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `SELECT * FROM product_history WHERE product_id=$1`
	rows, err := tx.Query(statement, product.Id)
	if err != nil {
		log.Printf("Error getting the product history of the product: %s. %s\n", product.Name, err.Error())
		return nil, err
	}
	defer rows.Close()

	var productHistory []models.ProductHistory
	for rows.Next() {
		product := &models.ProductHistory{}
		err = rows.Scan(&product.Id, &product.ProductId, &product.ProductName, &product.HigherPrice, &product.LowePrice, &product.OtherPaymentLowerPrice, &product.CreatedAt)
		if err != nil {
			log.Println("Error scanning the query result, getting the product history", err)
			return nil, err
		}
		productHistory = append(productHistory, *product)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return productHistory, nil
}
