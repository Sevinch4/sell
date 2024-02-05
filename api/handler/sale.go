package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
	"strconv"
)

// CreateSale godoc
// @Router       /sale [POST]
// @Summary      Create a new sale
// @Description  create a new sale
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 sale body models.CreateSale false "sale"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateSale(c *gin.Context) {
	sale := models.CreateSale{}
	if err := c.ShouldBindJSON(&sale); err != nil {
		handleResponse(c, "error is while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Sale().Create(context.Background(), sale)
	if err != nil {
		handleResponse(c, "error is while creating sale", http.StatusInternalServerError, err.Error())
		return
	}

	createdBranch, err := h.storage.Sale().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(
		c,
		"",
		http.StatusCreated,
		createdBranch,
	)
}

// GetSale godoc
// @Router       /sale/{id} [GET]
// @Summary      Get sale by id
// @Description  get sale by id
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 id path string true "sale_id"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetSale(c *gin.Context) {
	uid := c.Param("id")

	sale, err := h.storage.Sale().GetByID(context.Background(), uid)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, sale)
}

// GetSaleList godoc
// @Router       /sales [GET]
// @Summary      Get sale list
// @Description  get sale list
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetSaleList(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error is while converting page ", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error is while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	sales, err := h.storage.Sale().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	handleResponse(c, "", http.StatusOK, sales)
}

// UpdateSale godoc
// @Router       /sale/{id} [PUT]
// @Summary      Update sale
// @Description  update sale
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 id path string true "sale_id"
// @Param 		 sale body models.UpdateSale false "sale"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateSale(c *gin.Context) {
	uid := c.Param("id")
	sale := models.UpdateSale{}

	if err := c.ShouldBindJSON(&sale); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	sale.ID = uid
	id, err := h.storage.Sale().Update(context.Background(), sale)
	if err != nil {
		handleResponse(c, "error is while updating sale", http.StatusInternalServerError, err.Error())
		return
	}

	updatedSale, err := h.storage.Sale().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedSale)
}

// DeleteSale godoc
// @Router       /sale/{id} [DELETE]
// @Summary      Delete sale
// @Description  delete sale
// @Tags         sale
// @Accept       json
// @Produce      json
// @Param 		 id path string true "sale_id"
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteSale(c *gin.Context) {
	uid := c.Param("id")
	if err := h.storage.Sale().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error is while deleting", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "sale deleted!")
}
