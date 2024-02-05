package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"sell/api/models"
	"strconv"
)

// CreateTransaction godoc
// @Router       /transaction [POST]
// @Summary      Create a new transaction
// @Description  create a new transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param 		 transaction body models.CreateTransaction false "sale"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateTransaction(c *gin.Context) {
	trans := models.CreateTransaction{}
	if err := c.ShouldBindJSON(&trans); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Transaction().Create(context.Background(), trans)
	if err != nil {
		handleResponse(c, "error is while creating", http.StatusInternalServerError, err.Error())
		return
	}

	createdTrans, err := h.storage.Transaction().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdTrans)
}

// GetTransaction godoc
// @Router       /transaction/{id} [GET]
// @Summary      Get transaction by id
// @Description  get transaction by id
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param 		 id path string true "transaction_id"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetTransaction(c *gin.Context) {
	uid := c.Param("id")

	trans, err := h.storage.Transaction().GetByID(context.Background(), uid)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, trans)
}

// GetTransactionList godoc
// @Router       /transactions [GET]
// @Summary      Get transaction list
// @Description  get transaction list
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param		 page query string false "page"
// @Param		 limit query string false "limit"
// @Param		 from-amount query string false "from-amount"
// @Param		 to-amount query string false "to-amount"
// @Success      200  {object}  models.TransactionResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetTransactionList(c *gin.Context) {
	var (
		page, limit int
		fromAmount  float64
		toAmount    float64
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error is while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error is while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	fromAmountStr := c.DefaultQuery("from-amount", "0")
	fromAmount, err = strconv.ParseFloat(fromAmountStr, 64)
	if err != nil {
		handleResponse(c, "error is while converting from amount", http.StatusBadRequest, err.Error())
		return
	}

	toAmountStr := c.DefaultQuery("to-amount", fmt.Sprintf("%f", math.MaxFloat64))
	toAmount, err = strconv.ParseFloat(toAmountStr, 64)
	if err != nil {
		handleResponse(c, "error is while converting to amount", http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := h.storage.Transaction().GetList(context.Background(), models.TransactionGetListRequest{
		Page:       page,
		Limit:      limit,
		FromAmount: fromAmount,
		ToAmount:   toAmount,
	})

	if err != nil {
		handleResponse(c, "error is while getting list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, transactions)
}

// UpdateTransaction godoc
// @Router       /transaction/{id} [PUT]
// @Summary      Update transaction
// @Description  update transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param 		 id path string true "transaction_id"
// @Param 		 transaction body models.UpdateTransaction false "sale"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateTransaction(c *gin.Context) {
	uid := c.Param("id")

	trans := models.UpdateTransaction{}
	if err := c.ShouldBindJSON(&trans); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	trans.ID = uid

	id, err := h.storage.Transaction().Update(context.Background(), trans)
	if err != nil {
		handleResponse(c, "error is while updating trans", http.StatusInternalServerError, err.Error())
		return
	}

	updatedTrans, err := h.storage.Transaction().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedTrans)
}

// DeleteTransaction godoc
// @Router       /transaction/{id} [DELETE]
// @Summary      Delete transaction
// @Description  delete transaction
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Param 		 id path string true "transaction_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteTransaction(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.Transaction().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error is while deleting", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "transaction deleted!")
}
