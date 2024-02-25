package models

import "time"

type Staff struct {
	ID        string    `json:"id"`
	BranchID  string    `json:"branch_id"`
	TariffID  string    `json:"tariff_id"`
	StaffType string    `json:"staff_type"`
	Name      string    `json:"name"`
	Balance   uint      `json:"balance"`
	Age       uint      `json:"age"`
	BirthDate string    `json:"birth_date"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStaff struct {
	BranchID  string `json:"branch_id"`
	TariffID  string `json:"tariff_id"`
	StaffType string `json:"staff_type"`
	Name      string `json:"name"`
	Balance   uint   `json:"balance"`
	BirthDate string `json:"birth_date"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

type UpdateStaff struct {
	ID        string `json:"-"`
	BranchID  string `json:"branch_id"`
	TariffID  string `json:"tariff_id"`
	StaffType string `json:"staff_type"`
	Name      string `json:"name"`
	Balance   uint   `json:"balance"`
	Login     string `json:"login"`
}

type StaffsResponse struct {
	Staffs []Staff `json:"staffs"`
	Count  int     `json:"count"`
}

type UpdateStaffPassword struct {
	ID          string `json:"-"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type StaffType struct {
	ID      string `json:"id"`
	Balance uint   `json:"balance"`
}
type UpdateBalanceRequest struct {
	TransactionType string    `json:"transaction_type"`
	Source          string    `json:"source"`
	ShopAssistant   StaffType `json:"shop_assistant"`
	Cashier         StaffType `json:"cashier"`
	Text            string    `json:"text"`
	SaleID          string    `json:"sale_id"`
}
