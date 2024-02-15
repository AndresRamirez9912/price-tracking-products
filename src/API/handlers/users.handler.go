package handlers

import (
	"encoding/json"
	"net/http"
	apiModels "price-tracking-products/src/API/models"
	apiUtils "price-tracking-products/src/API/utils"
	"price-tracking-products/src/constants"
	"price-tracking-products/src/models"
	"price-tracking-products/src/services"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	response := &apiModels.GenericResponse{}
	userService, err := services.NewUserService()
	if err != nil {
		response.Success = false
		response.ErrorMessage = err.Error()
		apiUtils.CreateResponse(w, http.StatusInternalServerError, response)
		return
	}

	user := &models.User{}
	err = apiUtils.GetBody(r.Body, user)
	if err != nil {
		response.Success = false
		response.ErrorMessage = err.Error()
		apiUtils.CreateResponse(w, http.StatusBadRequest, response)
		return
	}

	err = userService.AddUser(*user)
	if err != nil {
		response.Success = false
		response.ErrorMessage = err.Error()
		apiUtils.CreateResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Send Response
	response.Success = true
	apiUtils.CreateResponse(w, http.StatusCreated, response)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userService, err := services.NewUserService()
	if err != nil {
		return
	}

	user := &apiModels.DeleteUserRequest{}
	err = apiUtils.GetBody(r.Body, user)
	if err != nil {
		return
	}

	err = userService.DeleteUser(user.User)
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusOK)
}

func ListUserProductsHandler(w http.ResponseWriter, r *http.Request) {
	response := &apiModels.ListUserProductsResponse{}
	userService, err := services.NewUserService()
	if err != nil {
		return
	}

	user := &apiModels.ListProductsRequest{}
	err = apiUtils.GetBody(r.Body, user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Success = false
		response.ErrorMessage = err.Error()
		return
	}

	products, err := userService.ListUserProducts(&user.User)
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
	w.Header().Add(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}
