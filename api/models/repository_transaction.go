package models

import "time"

type RepositoryTransaction struct {
	ID                        string     `json:"id"`
	ProductID                 string     `json:"product_id"`
	RepositoryTransactionType string     `json:"repository_transaction_type"`
	Price                     int        `json:"price"`
	Quantity                  int        `json:"quantity"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	DeletedAt                 *time.Time `json:"-"`
}

type CreateRepositoryTransaction struct {
	ProductID                 string `json:"product_id"`
	RepositoryTransactionType string `json:"repository_transaction_type"`
	Price                     int    `json:"price"`
	Quantity                  int    `json:"quantity"`
}

type UpdateRepositoryTransaction struct {
	ID                        string `json:"-"`
	ProductID                 string `json:"product_id"`
	RepositoryTransactionType string `json:"repository_transaction_type"`
	Price                     int    `json:"price"`
	Quantity                  int    `json:"quantity"`
}

type RepositoryTransactionsResponse struct {
	RepositoryTransactions []RepositoryTransaction `json:"repository_transactions"`
	Count                  int                     `json:"count"`
}
