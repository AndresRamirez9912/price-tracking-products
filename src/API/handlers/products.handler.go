package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	apiModels "price-tracking-products/src/API/models"
	apiUtils "price-tracking-products/src/API/utils"
)

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	body := &apiModels.AddProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		return
	}

	// Scrap the product information
	bodyRequest := &apiModels.ScrapingRequest{URL: body.URL}
	jsonData, err := json.Marshal(bodyRequest)
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:3000/scraping", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating the HTTP request to the Scraping service", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error sending the http request to the Scraping micro service", err)
		return
	}
	defer response.Body.Close()

	// Get the information
	product := &apiModels.ScrapProductResponse{}
	err = apiUtils.GetBody(response.Body, product)
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
