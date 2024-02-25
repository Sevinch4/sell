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
// @Param 		 status body models.SaleRequest true "status"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) EndSell(c *gin.Context) {
	saleID := c.Param("id")

	request := models.SaleRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	baskets, err := h.storage.Basket().GetList(context.Background(), models.GetListRequest{
		Page:   1,
		Limit:  10,
		Search: saleID,
	})
	if err != nil {
		handleResponse(c, "error is while getting  baskets list", http.StatusInternalServerError, err.Error())
		return
	}

	saleData, err := h.storage.Sale().GetByID(context.Background(), saleID)
	if err != nil {
		handleResponse(c, "error is while getting sale data", http.StatusInternalServerError, err.Error())
		return
	}
	saleTotalPrice := 0
	receivedProducts := make(map[string]models.Basket)

	for _, value := range baskets.Baskets {
		saleTotalPrice += value.Price
		receivedProducts[value.ProductID] = value
	}

	if request.Status == "cancel" {
		saleCancelID, err := h.storage.Sale().UpdatePrice(context.Background(), models.SaleRequest{
			SaleID:     saleID,
			TotalPrice: 0,
			Status:     "cancel",
		})
		if err != nil {
			handleResponse(c, "error is while updating cancel sale", http.StatusInternalServerError, err.Error())
			return
		}

		_, err = h.storage.Transaction().Create(context.Background(), models.CreateTransaction{
			SaleID:          saleID,
			StaffID:         saleData.ShopAssistantID,
			TransactionType: "withdraw",
			SourceType:      "sales",
			Amount:          0,
			Description:     "sale canceled",
		})

		handleResponse(c, "success", http.StatusOK, saleCancelID)
		return
	}

	updatedSalePrice, err := h.storage.Sale().UpdatePrice(context.Background(), models.SaleRequest{
		SaleID:     saleID,
		TotalPrice: saleTotalPrice,
		Status:     request.Status,
	})
	if err != nil {
		handleResponse(c, "error is while updating price", http.StatusInternalServerError, err.Error())
		return
	}

	response, err := h.storage.Sale().GetByID(context.Background(), updatedSalePrice)
	if err != nil {
		handleResponse(c, "error is while getting sale by updatedSalePrice", http.StatusInternalServerError, err.Error())
		return
	}

	repositoryData, err := h.storage.Repository().GetList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		handleResponse(c, "error while getting repositoryData list", http.StatusInternalServerError, err.Error())
		return
	}

	repoMap := make(map[string]models.Repository)
	for _, value := range repositoryData.Repositories {
		repoMap[value.ID] = value
	}

	for key, value := range repoMap {
		if value.ProductID == receivedProducts[value.ProductID].ProductID {
			_, err := h.storage.Repository().Update(context.Background(), models.UpdateRepository{
				ID:        key, // repoID
				ProductID: value.ProductID,
				BranchID:  value.BranchID,
				Count:     value.Count - receivedProducts[value.ProductID].Quantity, // repo_count - basket_quantity
			})
			if err != nil {
				handleResponse(c, "error while updating repositoryData prod quantities", http.StatusInternalServerError, err.Error())
				return
			}

			_, err = h.storage.RTransaction().Create(context.Background(), models.CreateRepositoryTransaction{
				ProductID:                 value.ProductID,
				RepositoryTransactionType: "minus",
				Price:                     receivedProducts[value.ProductID].Price,
				Quantity:                  receivedProducts[value.ProductID].Quantity,
			})
			if err != nil {
				handleResponse(c, "error while creating repositoryData transaction", http.StatusInternalServerError, err.Error())
				return
			}

		}
	}

	responseCashier, err := h.storage.Staff().StaffByID(context.Background(), models.PrimaryKey{ID: response.CashierID})
	if err != nil {
		handleResponse(c, "error while getting cashier by id", http.StatusInternalServerError, err.Error())
		return
	}
	responseCashierTariff, err := h.storage.StaffTariff().GetStaffTariffByID(context.Background(), models.PrimaryKey{ID: responseCashier.TariffID})
	if err != nil {
		handleResponse(c, "error while getting cashier tariff by id", http.StatusInternalServerError, err.Error())
		return
	}

	balance := 0

	if response.ShopAssistantID != "" {
		responseShopAssistant, err := h.storage.Staff().StaffByID(context.Background(), models.PrimaryKey{ID: response.ShopAssistantID})
		if err != nil {
			handleResponse(c, "error while getting shop assistant by id", http.StatusInternalServerError, err.Error())
			return
		}
		responseShopAssistantTariff, err := h.storage.StaffTariff().GetStaffTariffByID(context.Background(), models.PrimaryKey{ID: responseShopAssistant.TariffID})
		if err != nil {
			handleResponse(c, "error while getting shop assistant tariff by id", http.StatusInternalServerError, err.Error())
			return
		}
		if responseShopAssistantTariff.TariffType == "fixed" {
			if response.PaymentType == "cash" {
				balance += responseShopAssistantTariff.AmountForCash
			} else {
				balance += responseShopAssistantTariff.AmountForCard
			}
		} else if responseShopAssistantTariff.TariffType == "percent" {
			if response.PaymentType == "cash" {
				balance += (responseShopAssistantTariff.AmountForCash * saleTotalPrice) / 100
			} else {
				balance += (responseShopAssistantTariff.AmountForCard * saleTotalPrice) / 100
			}
		}

		//_, err = h.storage.Staff().UpdateStaff(context.Background(), models.UpdateStaff{
		//	ID:        responseShopAssistant.ID,
		//	BranchID:  response.BranchID,
		//	TariffID:  responseShopAssistant.TariffID,
		//	StaffType: responseShopAssistant.StaffType,
		//	Name:      responseShopAssistant.Name,
		//	Balance:   uint(balance),
		//	Login:     responseShopAssistant.Login,
		//})
		//if err != nil {
		//	handleResponse(c, "error is while updating staff", http.StatusInternalServerError, err.Error())
		//	return
		//}

		if responseCashierTariff.TariffType == "fixed" {
			if response.PaymentType == "cash" {
				balance += responseCashierTariff.AmountForCash
			} else {
				balance += responseCashierTariff.AmountForCard
			}
		} else if responseCashierTariff.TariffType == "percent" {
			if response.PaymentType == "cash" {
				balance += (responseCashierTariff.AmountForCash * saleTotalPrice) / 100
			} else {
				balance += (responseCashierTariff.AmountForCard * saleTotalPrice) / 100
			}
		}
		//_, err = h.storage.Staff().UpdateStaff(context.Background(), models.UpdateStaff{
		//	ID:        responseCashier.ID,
		//	BranchID:  response.BranchID,
		//	TariffID:  responseCashier.TariffID,
		//	StaffType: responseCashier.StaffType,
		//	Name:      responseCashier.Name,
		//	Balance:   uint(balance),
		//	Login:     responseCashier.Login,
		//})
		updateBalance := models.UpdateBalanceRequest{
			TransactionType: "topup",
			Source:          "sales",
			ShopAssistant: models.StaffType{
				ID:      responseShopAssistant.ID,
				Balance: uint(balance),
			},
			Cashier: models.StaffType{
				ID:      responseCashier.ID,
				Balance: uint(balance),
			},
			Text:   "some",
			SaleID: saleID,
		}

		if err := h.storage.Staff().UpdateBalance(context.Background(), updateBalance); err != nil {
			handleResponse(c, "error is while updating staff", http.StatusInternalServerError, err.Error())
			return
		}
	}

	handleResponse(c, "success", http.StatusOK, response)
}
