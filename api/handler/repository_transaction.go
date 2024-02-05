package handler

import (
	"context"
	"net/http"
	"sell/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRepositoryTransaction godoc
// @Router       /rtransaction [POST]
// @Summary      Create a new rtransaction
// @Description  create a new rtransaction
// @Tags         rtransaction
// @Accept       json
// @Produce      json
// @Param 		 rtransaction body models.CreateRepositoryTransaction false "rtransaction"
// @Success      200  {object}  models.RepositoryTransaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateRepositoryTransaction(c *gin.Context) {
	rtransaction := models.CreateRepositoryTransaction{}

	if err := c.ShouldBindJSON(&rtransaction); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.RTransaction().Create(context.Background(), rtransaction)
	if err != nil {
		handleResponse(c, "error while creating repository transaction", http.StatusInternalServerError, err.Error())
		return
	}

	createdRTransaction, err := h.storage.RTransaction().GetByID(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdRTransaction)
}

// GetRepositoryTransaction godoc
// @Router       /rtransaction/{id} [GET]
// @Summary      Get rtransaction by id
// @Description  get rtransaction by id
// @Tags         rtransaction
// @Accept       json
// @Produce      json
// @Param 		 id path string true "rtransaction_id"
// @Success      200  {object}  models.RepositoryTransaction
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetRepositoryTransaction(c *gin.Context) {
	uid := c.Param("id")

	repository, err := h.storage.RTransaction().GetByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting repository transaction by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, repository)
}

// GetRepositoryTransactionList godoc
// @Router       /rtransactions [GET]
// @Summary      Get rtransaction list
// @Description  get rtransaction list
// @Tags         rtransaction
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.RepositoryTransactionsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetRepositoryTransactionList(c *gin.Context) {
	var (
		page, limit int
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	search := c.Query("search")

	response, err := h.storage.RTransaction().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, "error while getting repository list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)
}

// UpdateRepositoryTransaction godoc
// @Router       /rtransaction/{id} [PUT]
// @Summary      Update rtransaction
// @Description  get rtransaction
// @Tags         rtransaction
// @Accept       json
// @Produce      json
// @Param 		 id path string true "rtransaction_id"
// @Param 		 rtransaction body models.UpdateRepositoryTransaction false "rtransaction"
// @Success      200  {object}  models.Repository
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateRepositoryTransaction(c *gin.Context) {
	uid := c.Param("id")

	rTransaction := models.UpdateRepositoryTransaction{}
	if err := c.ShouldBindJSON(&rTransaction); err != nil {
		handleResponse(c, "error while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	rTransaction.ID = uid
	if _, err := h.storage.RTransaction().Update(context.Background(), rTransaction); err != nil {
		handleResponse(c, "error while updating repository transaction ", http.StatusInternalServerError, err.Error())
		return
	}

	updatedRTransaction, err := h.storage.RTransaction().GetByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedRTransaction)
}

// DeleteRepositoryTransaction godoc
// @Router       /rtransaction/{id} [DELETE]
// @Summary      Delete rtransaction
// @Description  delete rtransaction
// @Tags         rtransaction
// @Accept       json
// @Produce      json
// @Param 		 id path string true "rtransaction_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteRepositoryTransaction(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.RTransaction().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error while deleting repository transaction ", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "repository transaction deleted")
}
