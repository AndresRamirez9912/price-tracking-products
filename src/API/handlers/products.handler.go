package handlers

import (
	"fmt"
	"net/http"
	apiModels "price-tracking-products/src/API/models"
	apiUtils "price-tracking-products/src/API/utils"
	"price-tracking-products/src/services"
)

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate: 1. If the user has the product, 2. If the product exists, 3. Scrap the product and link
	body := &apiModels.AddProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		return
	}

	// Check if the user already has the product
	hasProduct, err := body.User.HaveProduct(body.URL)
	if err != nil {
		return
	}

	if hasProduct {
		return
	}

	// Check if the product exists
	exists, err := body.User.ProductExists(body.URL)
	if err != nil {
		return
	}

	if exists {
		fmt.Println("The product exists")
		// Get the product By the URL
		product, err := body.User.GetProductByURL(body.URL)
		if err != nil {
			return
		}

		// Link Product with User in users_products DB
		err = product.AddProductToUser(body.User)
		if err != nil {
			return
		}
		fmt.Println("The product has added to the user")

		return
	}

	// Scrap the product information
	product, err := services.ScrapProduct(body.URL)
	if err != nil {
		return
	}

	// Add Product in Products DB
	err = product.Product.AddProduct()
	if err != nil {
		return
	}

	// Link Product with User in users_products DB
	err = product.Product.AddProductToUser(body.User)
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusCreated)
}

func RemoveProductHandler(w http.ResponseWriter, r *http.Request) {
	body := &apiModels.RemoveProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = body.Product.RemoveProductToUser(body.User)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send Response
	w.WriteHeader(http.StatusOK)
}
