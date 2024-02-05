package models

import "time"

type Basket struct {
	ID        string     `json:"id"`
	SaleID    string     `json:"sale_id"`
	ProductID string     `json:"product_id"`
	Quantity  int        `json:"quantity"`
	Price     int        `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type CreateBasket struct {
	SaleID    string `json:"sale_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}

type UpdateBasket struct {
	ID        string `json:"-"`
	SaleID    string `json:"sale_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}

type BasketsResponse struct {
	Baskets []Basket `json:"basket"`
	Count   int      `json:"count"`
}
