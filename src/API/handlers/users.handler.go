package handlers

import (
	"net/http"
	apiUtils "price-tracking-products/src/API/utils"
	"price-tracking-products/src/DB/repository"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &repository.User{}
	err := apiUtils.GetBody(r.Body, user)
	if err != nil {
		return
	}

	err = user.AddUser()
	if err != nil {
		return
	}

	// Send Response
	w.WriteHeader(http.StatusCreated)
}
