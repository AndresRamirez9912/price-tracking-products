package controller

import (
	"encoding/json"
	"net/http"

	"price-tracking-products/src/constants"
	apiModels "price-tracking-products/src/controller/models"
	apiUtils "price-tracking-products/src/controller/utils"
	"price-tracking-products/src/services"
)

type ProductsController struct {
	productService services.ProductServiceInterface
}

func NewProductController(productService services.ProductServiceInterface) *ProductsController {
	return &ProductsController{
		productService: productService,
	}
}

func (controller *ProductsController) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	body := &apiModels.AddProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		return
	}

	err = controller.productService.AddProduct(body.User, body.URL)
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusCreated)
}

func (controller *ProductsController) RemoveProductHandler(w http.ResponseWriter, r *http.Request) {
	body := &apiModels.RemoveProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = controller.productService.RemoveProductToUser(body.User, body.Product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send Response
	w.WriteHeader(http.StatusOK)
}

func (controller *ProductsController) GetProductHistory(w http.ResponseWriter, r *http.Request) {
	response := &apiModels.ProductHistoryResponse{}
	body := &apiModels.ProductHistoryRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	history, err := controller.productService.GetProductHistory(&body.Product)
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
