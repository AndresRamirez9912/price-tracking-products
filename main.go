package main

import (
	"log"
	"net/http"
	"price-tracking-products/src/API/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Post("/AddProduct", handlers.AddProductHandler)
	r.Post("/RemoveProduct", handlers.RemoveProductHandler)

	r.Post("/AddUser", handlers.AddUserHandler)

	log.Println("Starting server at port 3000")
	http.ListenAndServe(":3001", r)
}
