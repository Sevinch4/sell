package models

import "time"

type StaffTariff struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	TariffType    string     `json:"tariff_type"`
	AmountForCash int        `json:"amount_for_cash"`
	AmountForCard int        `json:"amount_for_card"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"-"`
}

type CreateStaffTariff struct {
	Name          string `json:"name"`
	TariffType    string `json:"tariff_type"`
	AmountForCash int    `json:"amount_for_cash"`
	AmountForCard int    `json:"amount_for_card"`
}

type UpdateStaffTariff struct {
	ID            string `json:"-"`
	Name          string `json:"name"`
	TariffType    string `json:"tariff_type"`
	AmountForCash int    `json:"amount_for_cash"`
	AmountForCard int    `json:"amount_for_card"`
}

type StaffTariffResponse struct {
	StaffTariffs []StaffTariff `json:"staff_tariffs"`
	Count        int           `json:"count"`
}
