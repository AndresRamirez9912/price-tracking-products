package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"price-tracking-products/src/constants"
	apiModels "price-tracking-products/src/controller/models"
	apiUtils "price-tracking-products/src/controller/utils"
	"price-tracking-products/src/models"
	"price-tracking-products/src/models/repository"
)

type ProductServiceInterface interface {
	AddProduct(user models.User, url string) error
	RemoveProductToUser(user models.User, product models.Product) error
	GetAllProducts() ([]models.Product, error)
	UpdateStoredProduct(product models.Product, available chan bool)

	GetProductHistory(product *models.Product) ([]models.ProductHistory, error)
}

type ProductService struct {
	repo     repository.ProductsRepository
	userRepo repository.UserRepository
}

func NewProductService(productRepo repository.ProductsRepository, userRepo repository.UserRepository) *ProductService {
	return &ProductService{repo: productRepo, userRepo: userRepo}
}

func (p *ProductService) AddProduct(user models.User, url string) error {
	// Check if the user already has the product
	hasProduct, err := p.userRepo.HaveProduct(user, url)
	if err != nil || hasProduct {
		return err
	}

	// Check if the product exists
	exists, err := p.repo.ProductExists(url)
	if err != nil {
		return err
	}

	if exists {
		// Get the product from the DB
		product, err := p.repo.GetProductByURL(url)
		if err != nil {
			return err
		}

		// Link the product to the user
		err = p.repo.AddProductToUser(user, *product)
		if err != nil {
			return err
		}
		return nil
	}

	// Scrap the product information
	scrapedProduct, err := scrapProduct(url)
	if err != nil {
		return err
	}

	// Add Product in Products DB
	err = p.repo.AddProduct(scrapedProduct.Product)
	if err != nil {
		return err
	}

	// Link Product with User in users_products DB
	err = p.repo.AddProductToUser(user, scrapedProduct.Product)
	if err != nil {
		return err
	}

	// Add product in history table
	err = p.repo.AddProductHistory(&scrapedProduct.Product)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductService) UpdateProductPrice(product models.Product) error {
	err := p.repo.UpdateProductPrice(product)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductService) RemoveProductToUser(user models.User, product models.Product) error {
	err := p.repo.RemoveProductToUser(user, product)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductService) GetProductHistory(product *models.Product) ([]models.ProductHistory, error) {
	history, err := p.repo.GetProductHistory(product)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (p *ProductService) GetAllProducts() ([]models.Product, error) {
	products, err := p.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductService) UpdateStoredProduct(product models.Product, available chan bool) {
	// Free an espace in the available cha (semaphore pattern)
	defer func() {
		<-available
	}()

	// Scrap the product
	newProduct, err := scrapProduct(product.ProductURL)
	if err != nil {
		log.Println("Error scrapping the product for price updating", err)
		return
	}

	// Compare if there is any change
	if !(newProduct.Product.Discount != product.Discount) {
		log.Println("Price equals")
		return
	}

	// Update the price
	err = p.repo.UpdateProductPrice(newProduct.Product)
	if err != nil {
		log.Println("Error updating the product price", err)
		return
	}

	// Add the price history
	err = p.repo.AddProductHistory(&newProduct.Product)
	if err != nil {
		log.Println("Error creating the item in the price history", err)
		return
	}

	// Notify
}

func scrapProduct(URL string) (*apiModels.ScrapProductResponse, error) {
	bodyRequest := &apiModels.ScrapingRequest{URL: URL}
	jsonData, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	// URL = "http://price-tracking-scrapping:3002/scraping"
	scrapingURL := fmt.Sprintf("%s://%s/api/scraping", os.Getenv(constants.SCHEME), os.Getenv(constants.SCRAPING_HOST))
	req, err := http.NewRequest(http.MethodPost, scrapingURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating the HTTP request to the Scraping service", err)
		return nil, err
	}
	req.Header.Add(constants.CONTENT_TYPE, constants.APPLICATION_JSON)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error sending the http request to the Scraping micro service", err)
		return nil, err
	}
	defer response.Body.Close()

	// Get the information
	product := &apiModels.ScrapProductResponse{}
	err = apiUtils.GetBody(response.Body, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}
