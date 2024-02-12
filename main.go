package main

import (
	"log"
	"net/http"
	"os"
	"price-tracking-products/src/API/handlers"
	"price-tracking-products/src/constants"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Read the .env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading the .env variables", err)
		return
	}

	r := chi.NewRouter()
	r.Post("/AddProduct", handlers.AddProductHandler)
	r.Post("/RemoveProduct", handlers.RemoveProductHandler)

	r.Post("/AddUser", handlers.AddUserHandler)
	r.Post("/DeleteUser", handlers.DeleteUserHandler)
	r.Post("/ListUserProducts", handlers.ListUserProductsHandler)

	log.Println("Starting server at port", os.Getenv(constants.PORT))
	http.ListenAndServe(os.Getenv(constants.PORT), r)
}
