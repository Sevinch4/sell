package storage

import (
	"context"
	"sell/api/models"
)

type IStorage interface {
	Close()
	StaffTariff() IStaffTariffRepo
	Staff() IStaffRepo
	Repository() IRepositoryRepo
	Basket() IBasketRepo
	RTransaction() IRepositoryTransactionRepo
	Category() ICategory
	Product() IProducts
	Branch() IBranchStorage
	Sale() ISaleStorage
	Transaction() ITransactionStorage
}

type IStaffTariffRepo interface {
	Create(context.Context, models.CreateStaffTariff) (string, error)
	GetStaffTariffByID(context.Context, models.PrimaryKey) (models.StaffTariff, error)
	GetStaffTariffList(context.Context, models.GetListRequest) (models.StaffTariffResponse, error)
	UpdateStaffTariff(context.Context, models.UpdateStaffTariff) (string, error)
	DeleteStaffTariff(context.Context, string) error
}

type IStaffRepo interface {
	Create(context.Context, models.CreateStaff) (string, error)
	StaffByID(context.Context, models.PrimaryKey) (models.Staff, error)
	GetStaffTList(context.Context, models.GetListRequest) (models.StaffsResponse, error)
	UpdateStaff(context.Context, models.UpdateStaff) (string, error)
	DeleteStaff(context.Context, string) error
	GetPassword(context.Context, string) (string, error)
	UpdatePassword(context.Context, models.UpdateStaffPassword) error
}

type IRepositoryRepo interface {
	Create(context.Context, models.CreateRepository) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.Repository, error)
	GetList(context.Context, models.GetListRequest) (models.RepositoriesResponse, error)
	Update(context.Context, models.UpdateRepository) (string, error)
	Delete(context.Context, string) error
}

type IBasketRepo interface {
	Create(context.Context, models.CreateBasket) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.Basket, error)
	GetList(context.Context, models.GetListRequest) (models.BasketsResponse, error)
	Update(context.Context, models.UpdateBasket) (string, error)
	Delete(context.Context, string) error
}

type IRepositoryTransactionRepo interface {
	Create(context.Context, models.CreateRepositoryTransaction) (string, error)
	GetByID(context.Context, models.PrimaryKey) (models.RepositoryTransaction, error)
	GetList(context.Context, models.GetListRequest) (models.RepositoryTransactionsResponse, error)
	Update(context.Context, models.UpdateRepositoryTransaction) (string, error)
	Delete(context.Context, string) error
}

type ICategory interface {
	Create(context.Context, models.CreateCategory) (string, error)
	GetByID(context.Context, string) (models.Category, error)
	GetList(context.Context, models.GetListRequest) (models.CategoryResponse, error)
	Update(context.Context, models.UpdateCategory) (string, error)
	Delete(context.Context, string) error
}

type IProducts interface {
	Create(context.Context, models.CreateProduct) (string, error)
	GetByID(context.Context, string) (models.Product, error)
	GetList(context.Context, models.ProductGetListRequest) (models.ProductResponse, error)
	Update(context.Context, models.UpdateProduct) (string, error)
	Delete(context.Context, string) error
}

type IBranchStorage interface {
	Create(context.Context, models.CreateBranch) (string, error)
	GetByID(context.Context, string) (models.Branch, error)
	GetList(context.Context, models.GetListRequest) (models.BranchResponse, error)
	Update(context.Context, models.UpdateBranch) (string, error)
	Delete(context.Context, string) error
}

type ISaleStorage interface {
	Create(context.Context, models.CreateSale) (string, error)
	GetByID(context.Context, string) (models.Sale, error)
	GetList(context.Context, models.GetListRequest) (models.SaleResponse, error)
	Update(context.Context, models.UpdateSale) (string, error)
	Delete(context.Context, string) error
	UpdatePrice(context.Context, int, string) (string, error)
}

type ITransactionStorage interface {
	Create(context.Context, models.CreateTransaction) (string, error)
	GetByID(context.Context, string) (models.Transaction, error)
	GetList(context.Context, models.TransactionGetListRequest) (models.TransactionResponse, error)
	Update(context.Context, models.UpdateTransaction) (string, error)
	Delete(context.Context, string) error
}
