package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
)

// Barcode godoc
// @Router       /barcode [POST]
// @Summary      barcode
// @Description  barcode
// @Tags         barcode
// @Accept       json
// @Produce      json
// @Param		 info body models.Barcode true "info"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) Barcode(c *gin.Context) {
	info := models.Barcode{}
	if err := c.ShouldBindJSON(&info); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	sale, err := h.storage.Sale().GetByID(context.Background(), info.SaleID)
	if err != nil {
		handleResponse(c, "error is getting sale by id", http.StatusInternalServerError, err.Error())
		return
	}

	if sale.Status == "success" {
		handleResponse(c, "sale ended", 300, "sale ended cannot add product")
		return
	}

	if sale.Status == "cancel" {
		handleResponse(c, "sale canceled", 300, "sale canceled cannot add product")
		return
	}

	products, err := h.storage.Product().GetList(context.Background(), models.ProductGetListRequest{
		Page:    1,
		Limit:   10,
		Barcode: info.Barcode,
	})

	if err != nil {
		handleResponse(c, "error is while getting product list by barcode", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		prodID    string
		prodPrice int
	)
	for _, product := range products.Products {
		prodID = product.ID
		prodPrice = product.Price
	}

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  10,
		Search: info.SaleID,
	})
	if err != nil {
		handleResponse(c, "error is while getting basket list", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		basketsMap = make(map[string]models.Basket)
		totalPrice = 0
	)

	// totalPrice
	totalPrice = info.Count * prodPrice

	for _, basket := range baskets.Baskets {
		basketsMap[basket.ProductID] = basket
	}

	repo, err := h.storage.Repository().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		handleResponse(c, "error is while getting repo list", http.StatusInternalServerError, err.Error())
		return
	}

	for _, r := range repo.Repositories {
		if prodID == basketsMap[r.ProductID].ProductID {
			// update un
			if r.Count < (basketsMap[r.ProductID].Quantity + info.Count) {
				handleResponse(c, "not enough product", 301, "not enough product")
				return
			}
		}

		// create un
		if r.Count < info.Count {
			handleResponse(c, "not enough product", 300, "not enough product")
			return
		}
	}

	isTrue := false

	for _, value := range basketsMap {
		if prodID == value.ProductID {
			isTrue = true
			id, err := h.storage.Basket().Update(context.Background(), models.UpdateBasket{
				ID:        value.ID,
				SaleID:    value.SaleID,
				ProductID: prodID,
				Quantity:  value.Quantity + info.Count,
				Price:     value.Price + totalPrice,
			})
			if err != nil {
				handleResponse(c, "error is while updating basket", 500, err.Error())
				return
			}
			updatedBasket, err := h.storage.Basket().GetByID(context.Background(), models.PrimaryKey{ID: id})
			if err != nil {
				handleResponse(c, "error is while getting basket", 500, err.Error())
				return
			}
			handleResponse(c, "updated", http.StatusOK, updatedBasket)
		}
	}

	if !isTrue {
		id, err := h.storage.Basket().Create(context.Background(), models.CreateBasket{
			SaleID:    info.SaleID,
			ProductID: prodID,
			Quantity:  info.Count,
			Price:     totalPrice,
		})
		if err != nil {
			handleResponse(c, "error is while creating basket", 500, err.Error())
			return
		}
		createdBasket, err := h.storage.Basket().GetByID(context.Background(), models.PrimaryKey{ID: id})
		if err != nil {
			handleResponse(c, "error is while getting basket", 500, err.Error())
			return
		}
		handleResponse(c, "updated", http.StatusOK, createdBasket)
	}
}
