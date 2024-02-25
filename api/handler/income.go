package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
	"strconv"
	"time"
)

// CreateIncome godoc
// @Router       /income [POST]
// @Summary      Create a new income
// @Description  create a new income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param 		 income body models.CreateIncome false "income"
// @Success      200  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateIncome(c *gin.Context) {
	income := models.CreateIncome{}
	if err := c.ShouldBindJSON(&income); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := h.storage.Income().Create(ctx, income)
	if err != nil {
		handleResponse(c, "error is while creating income", http.StatusInternalServerError, err.Error())
		return
	}

	createdIncome, err := h.storage.Income().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "success", http.StatusOK, createdIncome)
}

// GetIncome godoc
// @Router       /income/{id} [GET]
// @Summary      Get income
// @Description  get income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncome(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("id", id)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	income, err := h.storage.Income().GetByID(ctx, id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "success", http.StatusOK, income)
}

// GetIncomeList godoc
// @Router       /incomes [GET]
// @Summary      Get income list
// @Description  get income list
// @Tags         income
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 branch_id query string false "branch_id"
// @Success      200  {object}  models.IncomeResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetIncomeList(c *gin.Context) {
	var (
		page, limit int
		branchID    string
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

	branchID = c.Query("branch_id")
	resp, err := h.storage.Income().GetList(context.Background(), models.IncomeGetListRequest{
		Page:     page,
		Limit:    limit,
		BranchID: branchID,
	})
	if err != nil {
		handleResponse(c, "error is while getting incomes list", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(c, "success", http.StatusOK, resp)

}

// UpdateIncome godoc
// @Router       /income/{id} [PUT]
// @Summary      Update income
// @Description  update income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Param 		 income body models.UpdateIncome false "income"
// @Success      200  {object}  models.Income
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateIncome(c *gin.Context) {
	id := c.Param("id")

	income := models.UpdateIncome{}
	if err := c.ShouldBindJSON(&income); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	income.ID = id
	incomeID, err := h.storage.Income().Update(ctx, income)
	if err != nil {
		handleResponse(c, "error is while updating income", http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("icomeID", incomeID)

	updatedIncome, err := h.storage.Income().GetByID(ctx, incomeID)
	if err != nil {
		handleResponse(c, "error is while getting income", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "updated", http.StatusOK, updatedIncome)
}

// DeleteIncome godoc
// @Router       /income/{id} [DELETE]
// @Summary      Delete income
// @Description  delete income
// @Tags         income
// @Accept       json
// @Produce      json
// @Param 		 id path string true "income_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteIncome(c *gin.Context) {
	id := c.Param("id")
	if err := h.storage.Income().Delete(context.Background(), id); err != nil {
		handleResponse(c, "error is while deleting inocme", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "success", http.StatusOK, "income deleted!")
}
