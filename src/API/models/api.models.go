package apiModels

import (
	"price-tracking-products/src/models"
)

type AddProductRequest struct {
	URL  string      `json:"url"`
	User models.User `json:"user"`
}

type RemoveProductRequest struct {
	Product models.Product `json:"product"`
	User    models.User    `json:"user"`
}

type ListUserProductsResponse struct {
	Products []models.Product `json:"products"`
	GenericResponse
}

type ScrapingRequest struct {
	URL string `json:"url"`
}

type ScrapProductResponse struct {
	Product models.Product `json:"scrapedProduct"`
	GenericResponse
}

type GenericResponse struct {
	Success      bool   `json:"success"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
