package handler

import (
	"context"
	"net/http"
	"sell/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBasket godoc
// @Router       /basket [POST]
// @Summary      Create a new basket
// @Description  create a new basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 basket body models.CreateBasket false "basket"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateBasket(c *gin.Context) {
	basket := models.CreateBasket{}

	if err := c.ShouldBindJSON(&basket); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Basket().Create(context.Background(), basket)
	if err != nil {
		handleResponse(c, "error while creating basket", http.StatusInternalServerError, err.Error())
		return
	}

	createdBasket, err := h.storage.Basket().GetByID(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdBasket)
}

// GetBasket godoc
// @Router       /basket/{id} [GET]
// @Summary      Get basket by id
// @Description  get basket by id
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 id path string true "basket_id"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBasket(c *gin.Context) {
	uid := c.Param("id")

	basket, err := h.storage.Basket().GetByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting basket by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, basket)
}

// GetBasketList godoc
// @Router       /baskets [GET]
// @Summary      Get basket list
// @Description  get basket list
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.BasketsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBasketList(c *gin.Context) {
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

	response, err := h.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, "error while getting basket list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)
}

// UpdateBasket godoc
// @Router       /basket/{id} [PUT]
// @Summary      Update basket
// @Description  get basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 id path string true "basket_id"
// @Param 		 basket body models.UpdateBasket false "basket"
// @Success      200  {object}  models.Basket
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateBasket(c *gin.Context) {
	uid := c.Param("id")

	basket := models.UpdateBasket{}
	if err := c.ShouldBindJSON(&basket); err != nil {
		handleResponse(c, "error while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	basket.ID = uid
	if _, err := h.storage.Basket().Update(context.Background(), basket); err != nil {
		handleResponse(c, "error while updating basket ", http.StatusInternalServerError, err.Error())
		return
	}

	updatedBasket, err := h.storage.Basket().GetByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedBasket)
}

// DeleteBasket godoc
// @Router       /basket/{id} [DELETE]
// @Summary      Delete basket
// @Description  delete basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param 		 id path string true "basket_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteBasket(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.Basket().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error while deleting basket ", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "basket deleted")
}
