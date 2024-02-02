package apiModels

import "price-tracking-products/src/DB/repository"

type AddProductRequest struct {
	URL  string          `json:"url"`
	User repository.User `json:"user"`
}
