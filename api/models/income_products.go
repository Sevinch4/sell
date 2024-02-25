package models

type IncomeProduct struct {
	ID        string `json:"id"`
	IncomeID  string `json:"income_id"`
	ProductID string `json:"product_id"`
	Price     int    `json:"price"`
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateIncomeProduct struct {
	IncomeID  string `json:"income_id"`
	ProductID string `json:"product_id"`
	Price     int    `json:"price"`
	Count     int    `json:"count"`
}

type UpdateIncomeProduct struct {
	ID        string `json:"-"`
	IncomeID  string `json:"income_id"`
	ProductID string `json:"product_id"`
	Price     int    `json:"price"`
	Count     int    `json:"count"`
}

type IncomeProductsResponse struct {
	IncomeProducts        []IncomeProduct `json:"income_products"`
	CountOfIncomeProducts int             `json:"count_of_income_products"`
}

type IncomeProductRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	ProductID string `json:"product_id"`
	IncomeID  string `json:"income_id"`
}
