package main

import (
	"log"
	"net/http"
	"os"
	"price-tracking-products/src/constants"
	"price-tracking-products/src/controller"
	cronjob "price-tracking-products/src/controller/cron-job"
	"price-tracking-products/src/models/repository"
	"price-tracking-products/src/services"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize the application
	userRepo, err := repository.NewUserRepo()
	if err != nil {
		log.Fatal("Error initialyze the user repository", err)
	}
	productRepo, err := repository.NewProductRepo()
	if err != nil {
		log.Fatal("Error initialyze the product repository", err)
	}
	userService := services.NewUserService(userRepo)
	productService := services.NewProductService(productRepo, userRepo)
	userController := controller.NewUserController(userService)
	productController := controller.NewProductController(productService)

	// Init Cron Job
	go cronjob.InitCronJob(productService, 20) // 20 scraping at the same time

	r := chi.NewRouter()
	r.Post("/api/AddProduct", productController.AddProductHandler)
	r.Post("/api/RemoveProduct", productController.RemoveProductHandler)

	r.Post("/api/AddUser", userController.AddUserHandler)
	r.Post("/api/DeleteUser", userController.DeleteUserHandler)
	r.Post("/api/ListUserProducts", userController.ListUserProductsHandler)

	log.Println("Starting server at port", os.Getenv(constants.PORT))
	http.ListenAndServe(os.Getenv(constants.PORT), r)
}
