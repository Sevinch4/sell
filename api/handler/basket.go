package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
	"strconv"
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

	product, err := h.storage.Product().GetByID(context.Background(), basket.ProductID)
	if err != nil {
		handleResponse(c, "error is while getting product by id", http.StatusInternalServerError, err.Error())
	}

	totalSum := product.Price * basket.Quantity

	repo, err := h.storage.Repository().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  10,
		Search: basket.ProductID,
	})
	if err != nil {
		handleResponse(c, "error while getting repo", http.StatusInternalServerError, err.Error())
		return
	}

	for _, repository := range repo.Repositories {
		if repository.Count < basket.Quantity {
			handleResponse(c, "not enough product", 300, "not enough product")
			return
		}
	}

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  10,
		Search: basket.SaleID,
	})

	isTrue := false

	for _, value := range baskets.Baskets {
		if basket.ProductID == value.ProductID {
			for _, repository := range repo.Repositories {
				if repository.Count < basket.Quantity+value.Quantity {
					handleResponse(c, "not enough product", 300, "not enough product")
					return
				}
			}
			isTrue = true
			// Update
			id, err := h.storage.Basket().Update(context.Background(), models.UpdateBasket{
				ID:        value.ID,
				SaleID:    value.SaleID,
				ProductID: value.ProductID,
				Quantity:  value.Quantity + basket.Quantity,
				Price:     value.Price + totalSum,
			})
			if err != nil {
				handleResponse(c, "error is while updating basket", http.StatusInternalServerError, err.Error())
				return
			}
			updatedBasket, err := h.storage.Basket().GetByID(context.Background(), models.PrimaryKey{ID: id})
			if err != nil {
				handleResponse(c, "error is while getting basket by id", http.StatusInternalServerError, err.Error())
				return
			}
			handleResponse(c, "updated", http.StatusOK, updatedBasket)
		}

	}

	if !isTrue {
		//Create
		basket.Price = totalSum
		id, err := h.storage.Basket().Create(context.Background(), basket)
		if err != nil {
			handleResponse(c, "error while creating basket", http.StatusInternalServerError, err.Error())
			return
		}
		createdBasket, err := h.storage.Basket().GetByID(context.Background(), models.PrimaryKey{ID: id})
		if err != nil {
			handleResponse(c, "error is while getting basket by id", http.StatusInternalServerError, err.Error())
			return
		}
		handleResponse(c, "created", http.StatusOK, createdBasket)
	}

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
