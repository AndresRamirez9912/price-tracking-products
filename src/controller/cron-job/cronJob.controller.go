package cronjob

import (
	"price-tracking-products/src/services"
	"time"
)

func InitCronJob(productService services.ProductServiceInterface, workers int) {
	available := make(chan bool, workers)
	for {
		select {
		case <-time.After(1 * time.Hour):
			processProducts(productService, available)
		}
	}
}

func processProducts(productService services.ProductServiceInterface, available chan bool) {
	products, err := productService.GetAllProducts()
	if err != nil {
		return
	}

	// Use the semaphore pattern
	for _, product := range products {
		available <- true
		go productService.UpdateStoredProduct(product, available)
	}
}
