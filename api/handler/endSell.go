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

	totalPrice := 0
	receivedProducts := make(map[string]int)
	basketPrices := make(map[string]int)

	for _, value := range baskets.Baskets {
		totalPrice += value.Price
		receivedProducts[value.ProductID] = value.Quantity
		basketPrices[value.ProductID] = value.Price
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

	// storagedan prod quantity - quantity
	repo, err := h.storage.Repository().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		handleResponse(c, "error while getting repo list", http.StatusInternalServerError, err.Error())
		return
	}

	repoMap := make(map[string]models.Repository)
	for _, value := range repo.Repositories {
		repoMap[value.ID] = value
	}

	for key, value := range repoMap {
		_, err := h.storage.Repository().Update(context.Background(), models.UpdateRepository{
			ID:        key, // repoID
			ProductID: value.ProductID,
			BranchID:  value.BranchID,
			Count:     value.Count - receivedProducts[value.ProductID], // repo_count - basket_quantity
		})
		if err != nil {
			handleResponse(c, "error while updating repo prod quantities", http.StatusInternalServerError, err.Error())
			return
		}

		// storage_transaction -> check
		_, err = h.storage.RTransaction().Create(context.Background(), models.CreateRepositoryTransaction{
			StaffID:                   resp.CashierID,
			ProductID:                 value.ProductID,
			RepositoryTransactionType: "minus",
			Price:                     basketPrices[value.ProductID],
			Quantity:                  receivedProducts[value.ProductID],
		})
		if err != nil {
			handleResponse(c, "error while creating repo transaction", http.StatusInternalServerError, err.Error())
			return
		}
	}

	handleResponse(c, "success", http.StatusOK, resp)
}
