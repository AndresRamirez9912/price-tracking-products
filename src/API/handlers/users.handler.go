package handlers

import (
	"net/http"
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
