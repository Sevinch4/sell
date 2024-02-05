package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "sell/api/docs"
	"sell/api/handler"
	"sell/storage"
)

// New ...
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
func New(storage storage.IStorage) *gin.Engine {
	h := handler.New(storage)

	r := gin.New()

	r.POST("/sell", h.StartSell)
	r.PUT("/end-sell/:id", h.EndSell)

	r.POST("/category", h.CreateCategory)
	r.GET("/category/:id", h.GetCategory)
	r.GET("/categories", h.GetCategoryList)
	r.PUT("/category/:id", h.UpdateCategory)
	r.DELETE("/category/:id", h.DeleteCategory)

	r.POST("/product", h.CreateProduct)
	r.GET("/product/:id", h.GetProduct)
	r.GET("/products", h.GetProductList)
	r.PUT("/product/:id", h.UpdateProduct)
	r.DELETE("/product/:id", h.DeleteProduct)

	r.POST("/branch", h.CreateBranch)
	r.GET("/branch/:id", h.GetBranch)
	r.GET("/branches", h.GetBranchList)
	r.PUT("/branch/:id", h.UpdateBranch)
	r.DELETE("/branch/:id", h.DeleteBranch)

	r.POST("/repository", h.CreateRepository)
	r.GET("/repository/:id", h.GetRepository)
	r.GET("/repositories", h.GetRepositoryList)
	r.PUT("/repository/:id", h.UpdateRepository)
	r.DELETE("/repository/:id", h.DeleteRepository)

	r.POST("/sale", h.CreateSale)
	r.GET("/sale/:id", h.GetSale)
	r.GET("/sales", h.GetSaleList)
	r.PUT("/sale/:id", h.UpdateSale)
	r.DELETE("/sale/:id", h.DeleteSale)

	r.POST("/basket", h.CreateBasket)
	r.GET("/basket/:id", h.GetBasket)
	r.GET("/baskets", h.GetBasketList)
	r.PUT("/basket/:id", h.UpdateBasket)
	r.DELETE("/basket/:id", h.DeleteBasket)

	r.POST("/staff-tariff", h.CreateStaffTariff)
	r.GET("/staff-tariff/:id", h.GetStaffTariff)
	r.GET("/staff-tariffs", h.GetStaffTariffList)
	r.PUT("/staff-tariff/:id", h.UpdateStaffTariff)
	r.DELETE("/staff-tariff/:id", h.DeleteStaffTariff)

	r.POST("/staff", h.CreateStaff)
	r.GET("/staff/:id", h.GetStaff)
	r.GET("/staffs", h.GetStaffList)
	r.PUT("/staff/:id", h.UpdateStaff)
	r.PATCH("/staff/:id", h.UpdateStaffPassword)
	r.DELETE("/staff/:id", h.DeleteStaff)

	r.POST("/transaction", h.CreateTransaction)
	r.GET("/transaction/:id", h.GetTransaction)
	r.GET("/transactions", h.GetTransactionList)
	r.PUT("/transaction/:id", h.UpdateTransaction)
	r.DELETE("/transaction/:id", h.DeleteTransaction)

	r.POST("/rtransaction", h.CreateRepositoryTransaction)
	r.GET("/rtransaction/:id", h.GetRepositoryTransaction)
	r.GET("/rtransactions", h.GetRepositoryTransactionList)
	r.PUT("/rtransaction/:id", h.UpdateRepositoryTransaction)
	r.DELETE("/rtransaction/:id", h.DeleteRepositoryTransaction)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
	return r
}
