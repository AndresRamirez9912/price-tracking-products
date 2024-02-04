package apiModels

import "price-tracking-products/src/DB/repository"

type AddProductRequest struct {
	URL  string          `json:"url"`
	User repository.User `json:"user"`
}

type RemoveProductRequest struct {
	Product repository.Product `json:"product"`
	User    repository.User    `json:"user"`
}

type ScrapingRequest struct {
	URL string `json:"url"`
}

type ScrapProductResponse struct {
	Product      repository.Product `json:"scrapedProduct"`
	Success      bool               `json:"success"`
	ErrorCode    int                `json:"errorCode"`
	ErrorMessage string             `json:"errorMessage"`
}
