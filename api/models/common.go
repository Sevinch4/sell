package models

type PrimaryKey struct {
	ID string `json:"id"`
}

type GetListRequest struct {
	Page   int
	Limit  int
	Search string
}