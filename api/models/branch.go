package models

import (
	"time"
)

type Branch struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt string    `json:"-"`
}

type CreateBranch struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type UpdateBranch struct {
	ID      string `json:"-"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type BranchResponse struct {
	Branches []Branch
	Count    int
}
