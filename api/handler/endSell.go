package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
)

// EndSell godoc
// @Router       /end-sell/{id} [PUT]
// @Summary      end sell
// @Description  end sell
// @Tags         sell
// @Accept       json
// @Produce      json
// @Param 		 id path string true "sale_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) EndSell(c *gin.Context) {
	saleID := c.Param("id")

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  10,
		Search: saleID,
	})
	if err != nil {
		handleResponse(c, "error is while getting  baskets list", http.StatusInternalServerError, err.Error())
		return
	}

	var (
		productIDs        = make(map[string]string)
		productQuantities = make(map[string]int)
	)
	for _, value := range baskets.Baskets {
		productIDs[value.ID] = value.ProductID
		productQuantities[value.ID] = value.Quantity
	}

	products := make(map[string]int)

	for key, value := range productIDs {
		product, err := h.storage.Product().GetByID(context.Background(), value)
		if err != nil {
			handleResponse(c, "error is while getting product", http.StatusInternalServerError, err.Error())
			return
		}
		products[key] = product.Price
	}

	totalPrice := 0

	for key, value := range productQuantities {
		totalPrice += value * products[key]
	}

	id, err := h.storage.Sale().UpdatePrice(context.Background(), totalPrice, saleID)
	if err != nil {
		handleResponse(c, "error is while updating price", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.storage.Sale().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting sale by id", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(c, "success", http.StatusOK, resp)
}
