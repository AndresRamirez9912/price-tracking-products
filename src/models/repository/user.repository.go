package repository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"price-tracking-products/src/constants"
	apiUtils "price-tracking-products/src/controller/utils"
	"price-tracking-products/src/models"
	dbUtils "price-tracking-products/src/models/utils"
)

type UserRepository interface {
	AddUser(models.User) error
	DeleteUser(models.User) error
	ListUserProducts(models.User) ([]models.Product, error)
	HaveProduct(models.User, string) (bool, error)
	DeleteAllUserProducts(models.User) error
	CreateUser(user models.User) (*models.CreateUserResponse, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo() (*UserRepo, error) {
	db, err := dbUtils.OpenDBConnection()
	if err != nil {
		return nil, err
	}
	return &UserRepo{db: db}, nil
}

func (repo UserRepo) AddUser(user models.User) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `insert into "users" ("id","email","user_name") values ($1, $2, $3)`
	_, err = tx.Exec(statement, user.Id, user.Email, user.UserName)
	if err != nil {
		log.Printf("Error adding %s user %s\n", user.UserName, err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo UserRepo) DeleteUser(user models.User) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `DELETE FROM "users" WHERE id=$1`
	_, err = tx.Exec(statement, user.Id)
	if err != nil {
		log.Printf("Error deleting the user %s. %s \n", user.UserName, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo UserRepo) ListUserProducts(user models.User) ([]models.Product, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if (user == models.User{}) {
		log.Println("User empty, bad request")
		return nil, errors.New("User empty, bad request")
	}

	statement := `SELECT id,name,brand,higher_price,lower_price,other_price,discount,image_url,product_url,store from "products" 
	INNER JOIN "products_users" ON products_users.product_id = products.id 
	WHERE products_users.user_id = $1`
	rows, err := tx.Query(statement, user.Id)
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

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (repo UserRepo) HaveProduct(user models.User, url string) (bool, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `SELECT count(user_id) FROM products_users 
	INNER JOIN products ON products_users.product_id = products.id
	WHERE user_id=$1 AND product_url=$2;`
	rows, err := tx.Query(statement, user.Id, url)
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

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return false, nil
}

func (repo UserRepo) DeleteAllUserProducts(user models.User) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statement := `DELETE FROM "products_users" WHERE user_id=$1`
	_, err = tx.Exec(statement, user.Id)
	if err != nil {
		log.Printf("Error deleting the user's products. %s \n", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (repo UserRepo) CreateUser(user models.User) (*models.CreateUserResponse, error) {
	bodyRequest := &models.CreateUserRequest{
		Name:     user.UserName,
		Email:    user.Email,
		UserName: user.UserName,
		Password: user.Password,
	}
	jsonData, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	// URL = "http://price-tracking-auth:3001/signUp"
	scrapingURL := fmt.Sprintf("%s://%s/api/signUp", os.Getenv(constants.SCHEME), os.Getenv(constants.AUTH_HOST))
	req, err := http.NewRequest(http.MethodPost, scrapingURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating the HTTP request to the Auth service", err)
		return nil, err
	}
	req.Header.Add(constants.CONTENT_TYPE, constants.APPLICATION_JSON)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error sending the http request to the Auth micro service", err)
		return nil, err
	}
	defer response.Body.Close()

	// Get the information
	product := &models.CreateUserResponse{}
	err = apiUtils.GetBody(response.Body, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}
