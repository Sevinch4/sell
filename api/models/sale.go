package models

import "time"

type Sale struct {
	ID              string    `json:"id"`
	BranchID        string    `json:"branch_id"`
	ShopAssistantID string    `json:"shop_assistant_id"`
	CashierID       string    `json:"cashier_id"`
	PaymentType     string    `json:"payment_type"`
	Price           float32   `json:"price"`
	Status          string    `json:"status"`
	ClientName      string    `json:"client_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateSale struct {
	BranchID        string  `json:"branch_id"`
	ShopAssistantID string  `json:"shop_assistant_id"`
	CashierID       string  `json:"cashier_id"`
	PaymentType     string  `json:"payment_type"`
	Price           float32 `json:"price"`
	Status          string  `json:"status"`
	ClientName      string  `json:"client_name"`
}

type UpdateSale struct {
	ID              string  `json:"-"`
	BranchID        string  `json:"branch_id"`
	ShopAssistantID string  `json:"shop_assistant_id"`
	CashierID       string  `json:"cashier_id"`
	PaymentType     string  `json:"payment_type"`
	Price           float32 `json:"price"`
	Status          string  `json:"status"`
	ClientName      string  `json:"client_name"`
}

type SaleResponse struct {
	Sales []Sale
	Count int
}
