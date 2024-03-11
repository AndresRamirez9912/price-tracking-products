package controller

import (
	"encoding/json"
	"net/http"

	"price-tracking-products/src/constants"
	apiModels "price-tracking-products/src/controller/models"
	apiUtils "price-tracking-products/src/controller/utils"
	"price-tracking-products/src/models"
	"price-tracking-products/src/services"
)

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(userService services.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (controller *UserController) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	response := &apiModels.GenericResponse{}
	user := &models.User{}
	err := apiUtils.GetBody(r.Body, user)
	if err != nil {
		response.Success = false
		response.ErrorMessage = err.Error()
		apiUtils.CreateResponse(w, http.StatusBadRequest, response)
		return
	}

	err = controller.userService.AddUser(*user)
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

func (controller *UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &apiModels.DeleteUserRequest{}
	err := apiUtils.GetBody(r.Body, user)
	if err != nil {
		return
	}

	err = controller.userService.DeleteUser(user.User)
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusOK)
}

func (controller *UserController) ListUserProductsHandler(w http.ResponseWriter, r *http.Request) {
	response := &apiModels.ListUserProductsResponse{}
	user := &apiModels.ListProductsRequest{}
	err := apiUtils.GetBody(r.Body, user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Success = false
		response.ErrorMessage = err.Error()
		return
	}

	products, err := controller.userService.ListUserProducts(&user.User)
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
