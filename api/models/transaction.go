package models

import (
	"time"
)

type Transaction struct {
	ID              string    `json:"id"`
	SaleID          string    `json:"sale_id"`
	StaffID         string    `json:"staff_id"`
	TransactionType string    `json:"transaction_type"`
	SourceType      string    `json:"source_type"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       string    `json:"-"`
}

type CreateTransaction struct {
	SaleID          string  `json:"sale_id"`
	StaffID         string  `json:"staff_id"`
	TransactionType string  `json:"transaction_type"`
	SourceType      string  `json:"source_type"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
}

type UpdateTransaction struct {
	ID              string  `json:"-"`
	SaleID          string  `json:"sale_id"`
	StaffID         string  `json:"staff_id"`
	TransactionType string  `json:"transaction_type"`
	SourceType      string  `json:"source_type"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
}

type TransactionResponse struct {
	Transactions []Transaction
	Count        int
}

type TransactionGetListRequest struct {
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	FromAmount float64 `json:"from_amount"`
	ToAmount   float64 `json:"to_amount"`
}
