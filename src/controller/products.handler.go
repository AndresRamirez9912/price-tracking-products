package controller

import (
	"net/http"

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
		apiUtils.CreateErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = controller.productService.AddProduct(body.User, body.URL)
	if err != nil {
		apiUtils.CreateErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send Response
	w.WriteHeader(http.StatusCreated)
}

func (controller *ProductsController) RemoveProductHandler(w http.ResponseWriter, r *http.Request) {
	body := &apiModels.RemoveProductRequest{}
	err := apiUtils.GetBody(r.Body, body)
	if err != nil {
		apiUtils.CreateErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = controller.productService.RemoveProductToUser(body.User, body.Product)
	if err != nil {
		apiUtils.CreateErrorResponse(w, http.StatusInternalServerError, err.Error())
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
		apiUtils.CreateErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	history, err := controller.productService.GetProductHistory(&body.Product)
	if err != nil {
		apiUtils.CreateErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send Response
	response.History = history
	response.Success = true
	apiUtils.CreateResponse(w, http.StatusOK, response)
}
