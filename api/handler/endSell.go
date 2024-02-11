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

	totalPrice := 0
	receivedProducts := make(map[string]models.Basket)

	for _, value := range baskets.Baskets {
		totalPrice += value.Price
		receivedProducts[value.ProductID] = value
	}

	if request.Status == "cancel" {
		totalPrice = 0
	}

	updatedSalePrice, err := h.storage.Sale().UpdatePrice(context.Background(), models.SaleRequest{
		SaleID:     saleID,
		TotalPrice: totalPrice,
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

	if request.Status == "cancel" {
		_, err = h.storage.Transaction().Create(context.Background(), models.CreateTransaction{
			SaleID:          saleID,
			StaffID:         response.CashierID,
			TransactionType: "withdraw",
			SourceType:      "sales",
			Amount:          float64(totalPrice),
			Description:     "sale canceled nothing sold",
		})
		if err != nil {
			handleResponse(c, "error is while creating transaction for staff", http.StatusInternalServerError, err.Error())
			return
		}
		handleResponse(c, "success", http.StatusOK, response)
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
				StaffID:                   response.CashierID,
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

	tariffs, err := h.storage.StaffTariff().GetStaffTariffList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		handleResponse(c, "error is while getting staff tariff list", http.StatusBadRequest, err.Error())
		return
	}

	staffs, err := h.storage.Staff().GetStaffTList(context.Background(), models.GetListRequest{
		Page:  1,
		Limit: 10,
	})

	staffsMap := make(map[string]models.Staff)
	staffBalance := make(map[string]uint)
	for _, staff := range staffs.Staffs {
		if response.CashierID == staff.ID || response.ShopAssistantID == staff.ID {
			staffsMap[staff.TariffID] = staff
			staffBalance[staff.TariffID] = staff.Balance
		}
	}

	tariffsMap := make(map[string]models.StaffTariff)
	for _, tariff := range tariffs.StaffTariffs {
		if tariffsMap[tariff.ID].ID == staffsMap[tariff.ID].TariffID {
			continue
		}
		tariffsMap[tariff.ID] = tariff
	}

	for _, tariff := range tariffsMap {
		switch {
		case tariff.TariffType == "fixed":
			switch response.PaymentType {
			case "cash":
				staffBalance[tariff.ID] += uint(tariff.AmountForCash)
			case "card":
				staffBalance[tariff.ID] += uint(tariff.AmountForCard)
			}
		case tariff.TariffType == "percent":
			switch response.PaymentType {
			case "cash":
				staffBalance[tariff.ID] += uint(totalPrice * tariff.AmountForCash / 100)
			case "card":
				staffBalance[tariff.ID] += uint(totalPrice * tariff.AmountForCard / 100)
			}
		}

		_, err := h.storage.Staff().UpdateStaff(context.Background(), models.UpdateStaff{
			ID:        staffsMap[tariff.ID].ID,
			BranchID:  staffsMap[tariff.ID].BranchID,
			TariffID:  staffsMap[tariff.ID].TariffID,
			StaffType: staffsMap[tariff.ID].StaffType,
			Name:      staffsMap[tariff.ID].Name,
			Balance:   staffBalance[tariff.ID],
			Login:     staffsMap[tariff.ID].Login,
		})
		if err != nil {
			handleResponse(c, "error is while updating staff", http.StatusInternalServerError, err.Error())
			return
		}

		_, err = h.storage.Transaction().Create(context.Background(), models.CreateTransaction{
			SaleID:          saleID,
			StaffID:         staffsMap[tariff.ID].ID,
			TransactionType: "topup",
			SourceType:      "sales",
			Amount:          float64(totalPrice),
			Description:     "staff sell products",
		})
		if err != nil {
			handleResponse(c, "error is while creating transaction for staff", http.StatusInternalServerError, err.Error())
			return
		}
	}

	handleResponse(c, "success", http.StatusOK, response)
}
