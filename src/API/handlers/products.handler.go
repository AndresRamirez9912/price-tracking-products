package handlers

import (
	"encoding/json"
	"net/http"
	apiModels "price-tracking-products/src/API/models"
	apiUtils "price-tracking-products/src/API/utils"
	"price-tracking-products/src/constants"
	"price-tracking-products/src/services"
)

func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	productService := services.NewProductService()
	body := &apiModels.AddProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		return
	}

	err = productService.AddProduct(body.User, body.URL)
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusCreated)
}

func RemoveProductHandler(w http.ResponseWriter, r *http.Request) {
	productService := services.NewProductService()
	body := &apiModels.RemoveProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = productService.RemoveProductToUser(body.User, body.Product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send Response
	w.WriteHeader(http.StatusOK)
}

func GetProductHistory(w http.ResponseWriter, r *http.Request) {
	response := &apiModels.ProductHistoryResponse{}
	productService := services.NewProductService()
	body := &apiModels.ProductHistoryRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	history, err := productService.GetProductHistory(&body.Product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create response
	response.History = history
	response.Success = true
	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Success = false
		response.ErrorMessage = err.Error()
		return
	}

	// Send Response
	w.Header().Add(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
