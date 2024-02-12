package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	apiModels "price-tracking-products/src/API/models"
	apiUtils "price-tracking-products/src/API/utils"
	"price-tracking-products/src/DB/repository"
	"price-tracking-products/src/models"
)

type ProductService struct {
	repo repository.ProductsRepository
}

func NewProductService() *ProductService {
	repo := repository.NewProductRepo()
	return &ProductService{repo: repo}
}

func (p *ProductService) AddProduct(user models.User, url string) error {
	// Check if the user already has the product
	userService, err := NewUserService()
	if err != nil {
		return err
	}
	hasProduct, err := userService.repo.HaveProduct(user, url)
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
	scrapedProduct, err := ScrapProduct(url)
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
	return nil
}

func (p *ProductService) DeleteProduct(product models.Product) error {
	err := p.repo.DeleteProduct(product)
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

func (p *ProductService) GetProductByURL(url string) (*models.Product, error) {
	product, err := p.repo.GetProductByURL(url)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductService) CheckIfProductExists(url string) (bool, error) {
	exists, err := p.repo.ProductExists(url)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *ProductService) AddProductToUser(user models.User, product models.Product) error {
	err := p.repo.AddProductToUser(user, product)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductService) RemoveProductToUser(user models.User, product models.Product) error {
	err := p.repo.AddProductToUser(user, product)
	if err != nil {
		return err
	}
	return nil
}

func ScrapProduct(URL string) (*apiModels.ScrapProductResponse, error) {
	bodyRequest := &apiModels.ScrapingRequest{URL: URL}
	jsonData, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "http://price-tracking-scrapping:3002/scraping", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating the HTTP request to the Scraping service", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

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
