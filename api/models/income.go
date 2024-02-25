package models

import "time"

type Income struct {
	ID        string    `json:"id"`
	BranchID  string    `json:"branch_id"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateIncome struct {
	BranchID string `json:"branch_id"`
}

type UpdateIncome struct {
	ID       string `json:"-"`
	BranchID string `json:"branch_id"`
	Price    int    `json:"price"`
}

type IncomeResponse struct {
	Incomes []Income `json:"incomes"`
	Count   int      `json:"count"`
}

type IncomeGetListRequest struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	BranchID string `json:"branch_id"`
}
