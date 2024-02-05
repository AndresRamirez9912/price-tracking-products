package handlers

import (
	"encoding/json"
	"net/http"
	apiModels "price-tracking-products/src/API/models"
	apiUtils "price-tracking-products/src/API/utils"
	"price-tracking-products/src/models"
	"price-tracking-products/src/services"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	userService := services.NewUserService()
	user := &models.User{}
	err := apiUtils.GetBody(r.Body, user)
	if err != nil {
		return
	}

	err = userService.AddUser(*user)
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusCreated)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userService := services.NewUserService()
	user := &models.User{}
	err := apiUtils.GetBody(r.Body, user)
	if err != nil {
		return
	}

	err = userService.DeleteUser(*user)
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusCreated)
}

func ListUserProductsHandler(w http.ResponseWriter, r *http.Request) {
	response := &apiModels.ListUserProductsResponse{}
	userService := services.NewUserService()
	user := &models.User{}
	err := apiUtils.GetBody(r.Body, user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Success = false
		response.ErrorMessage = err.Error()
		return
	}

	products, err := userService.ListUserProducts(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Success = false
		response.ErrorMessage = err.Error()
		return
	}

	response.Products = products
	response.Success = true
	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Success = false
		response.ErrorMessage = err.Error()
		return
	}

	// Send Response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}
