package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
	"strconv"
	"time"
)

// CreateIncomeProduct godoc
// @Router       /income-product [POST]
// @Summary      Create a new income products
// @Description  create a new income products
// @Tags         income-product
// @Accept       json
// @Produce      json
// @Param 		 income-product body models.CreateIncomeProduct false "income-product"
// @Success      200  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateIncomeProduct(c *gin.Context) {
	incomeProduct := models.CreateIncomeProduct{}
	if err := c.ShouldBindJSON(&incomeProduct); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := h.storage.IncomeProducts().Create(ctx, incomeProduct)
	if err != nil {
		handleResponse(c, "error is while creating income Product", http.StatusInternalServerError, err.Error())
		return
	}

	incomeProductsList, err := h.storage.IncomeProducts().GetList(ctx, models.IncomeProductRequest{
		Page:     1,
		Limit:    10,
		IncomeID: incomeProduct.IncomeID,
	})
	if err != nil {
		handleResponse(c, "error is while getting list income products", http.StatusInternalServerError, err.Error())
		return
	}

	incomePrice := 0
	for _, i := range incomeProductsList.IncomeProducts {
		incomePrice += i.Price
	}

	repo, err := h.storage.Repository().GetList(ctx, models.GetListRequest{
		Page:   1,
		Limit:  10,
		Search: incomeProduct.ProductID,
	})
	if err != nil {
		handleResponse(c, "error is while getting repo list", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		repoID       string
		repoBranchID string
		inRepoCount  int
	)

	for _, r := range repo.Repositories {
		repoID = r.ID
		repoBranchID = r.BranchID
		inRepoCount = r.Count
	}

	_, err = h.storage.Repository().Update(ctx, models.UpdateRepository{
		ID:        repoID,
		ProductID: incomeProduct.ProductID,
		BranchID:  repoBranchID,
		Count:     inRepoCount + incomeProduct.Count,
	})
	if err != nil {
		handleResponse(c, "error is while updating repo", http.StatusInternalServerError, err.Error())
		return
	}

	_, err = h.storage.Income().Update(ctx, models.UpdateIncome{
		ID:       incomeProduct.IncomeID,
		BranchID: repoBranchID,
		Price:    incomePrice,
	})
	if err != nil {
		handleResponse(c, "error is while updating income", http.StatusInternalServerError, err.Error())
		return
	}

	for _, value := range incomeProductsList.IncomeProducts {
		//fmt.Println("staffID", staffID[value.ProductID])
		_, err := h.storage.RTransaction().Create(ctx, models.CreateRepositoryTransaction{
			ProductID:                 value.ProductID,
			RepositoryTransactionType: "plus",
			Price:                     value.Price,
			Quantity:                  value.Count,
		})
		if err != nil {
			handleResponse(c, "error is while creating repo transaction", http.StatusInternalServerError, err.Error())
			return
		}
	}

	createdIncomeProduct, err := h.storage.IncomeProducts().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "success", http.StatusOK, createdIncomeProduct)
}

// GetIncomeProduct godoc
// @Router       /income-product/{id} [GET]
// @Summary      Get income product
// @Description  get income product
// @Tags         income-product
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeProduct(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	incomeProducts, err := h.storage.IncomeProducts().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "success", http.StatusOK, incomeProducts)
}

// GetIncomeProductsList godoc
// @Router       /income-product [GET]
// @Summary      Get income product list
// @Description  get income product list
// @Tags         income-product
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 branch_id query string false "branch_id"
// @Success      200  {object}  models.IncomeProductsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeProductsList(c *gin.Context) {
	var (
		page, limit int
		productID   string
		err         error
	)
	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error is while converting pageStr", http.StatusBadRequest, err)
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error is while converting limitStr", http.StatusBadRequest, err)
		return
	}

	productID = c.Query("product_id")
	resp, err := h.storage.IncomeProducts().GetList(context.Background(), models.IncomeProductRequest{
		Page:      page,
		Limit:     limit,
		ProductID: productID,
	})
	if err != nil {
		handleResponse(c, "error is while getting income product list", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(c, "success", http.StatusOK, resp)

}

// UpdateIncomeProduct godoc
// @Router       /income-product/{id} [PUT]
// @Summary      Update income product
// @Description  update income product
// @Tags         income-product
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Param 		 income-product body models.UpdateIncomeProduct false "income-product"
// @Success      200  {object}  models.IncomeProduct
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateIncomeProduct(c *gin.Context) {
	id := c.Param("id")

	incomeProduct := models.UpdateIncomeProduct{}
	if err := c.ShouldBindJSON(&incomeProduct); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	incomeProduct.ID = id
	incomeID, err := h.storage.IncomeProducts().Update(ctx, incomeProduct)
	if err != nil {
		handleResponse(c, "error is while updating income", http.StatusInternalServerError, err.Error())
		return
	}

	updatedIncomeProduct, err := h.storage.IncomeProducts().GetByID(ctx, incomeID)
	if err != nil {
		handleResponse(c, "error is while getting income", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "updated", http.StatusOK, updatedIncomeProduct)
}

// DeleteIncomeProduct godoc
// @Router       /income-product/{id} [DELETE]
// @Summary      Delete income product
// @Description  delete income product
// @Tags         income-product
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteIncomeProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.storage.IncomeProducts().Delete(context.Background(), id); err != nil {
		handleResponse(c, "error is while deleting income", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "success", http.StatusOK, "income product deleted!")
}
