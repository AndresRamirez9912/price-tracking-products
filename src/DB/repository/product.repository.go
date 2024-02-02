package repository

type Product struct {
	Name                   string `json:"name"`
	Id                     string `json:"id"`
	Brand                  string `json:"brand"`
	HigherPrice            string `json:"higherPrice"`
	LowePrice              string `json:"lowePrice"`
	OtherPaymentLowerPrice string `json:"otherPaymentLowerPrice"`
	Discount               string `json:"discount"`
	ImageURL               string `json:"imageURL"`
	Store                  string `json:"store"`
	ProductURL             string `json:"productURL"`
}

type ProductsRepository interface {
	AddProduct(Product) (Product, error)
	DeleteProduct(Product) error
}
