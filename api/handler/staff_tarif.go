package handler

import (
	"context"
	"net/http"
	"sell/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateStaffTariff godoc
// @Router       /staff-tariff [POST]
// @Summary      Create a new staff tariff
// @Description  create a new staff tariff
// @Tags         staff-tariff
// @Accept       json
// @Produce      json
// @Param 		 staffTariff body models.CreateStaffTariff false "staff-Tariff"
// @Success      200  {object}  models.StaffTariff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateStaffTariff(c *gin.Context) {
	staffTariff := models.CreateStaffTariff{}

	if err := c.ShouldBindJSON(&staffTariff); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.StaffTariff().Create(context.Background(), staffTariff)
	if err != nil {
		handleResponse(c, "error while creating staff tariff", http.StatusInternalServerError, err.Error())
		return
	}

	createdStaffTariff, err := h.storage.StaffTariff().GetStaffTariffByID(context.Background(), models.PrimaryKey{ID: id})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdStaffTariff)
}

// GetStaffTariff godoc
// @Router       /staff-tariff/{id} [GET]
// @Summary      Get staff tariff by id
// @Description  get staff tariff by id
// @Tags         staff-tariff
// @Accept       json
// @Produce      json
// @Param 		 id path string true "staff-tariff_id"
// @Success      200  {object}  models.StaffTariff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStaffTariff(c *gin.Context) {
	uid := c.Param("id")

	staffTariff, err := h.storage.StaffTariff().GetStaffTariffByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting staff tariff by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, staffTariff)
}

// GetStaffTariffList godoc
// @Router       /staff-tariffs [GET]
// @Summary      Get staff tariff list
// @Description  get staff tariff list
// @Tags         staff-tariff
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.StaffTariffResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStaffTariffList(c *gin.Context) {
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

	response, err := h.storage.StaffTariff().GetStaffTariffList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, "error while getting staff tariff list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)
}

// UpdateStaffTariff godoc
// @Router       /staff-tariff/{id} [PUT]
// @Summary      Update staff tariff
// @Description  get staff tariff
// @Tags         staff-tariff
// @Accept       json
// @Produce      json
// @Param 		 id path string true "staff-tariff_id"
// @Param 		 staff-tariff body models.UpdateStaffTariff false "staff-tariff"
// @Success      200  {object}  models.StaffTariff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateStaffTariff(c *gin.Context) {
	uid := c.Param("id")

	sTariff := models.UpdateStaffTariff{}
	if err := c.ShouldBindJSON(&sTariff); err != nil {
		handleResponse(c, "error while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	sTariff.ID = uid
	if _, err := h.storage.StaffTariff().UpdateStaffTariff(context.Background(), sTariff); err != nil {
		handleResponse(c, "error while updating staff tariff", http.StatusInternalServerError, err.Error())
		return
	}

	updatedStaffTariff, err := h.storage.StaffTariff().GetStaffTariffByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedStaffTariff)
}

// DeleteStaffTariff godoc
// @Router       /staff-tariff/{id} [DELETE]
// @Summary      Delete staff tariff
// @Description  delete staff tariff
// @Tags         staff-tariff
// @Accept       json
// @Produce      json
// @Param 		 id path string true "staff-tariff_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteStaffTariff(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.StaffTariff().DeleteStaffTariff(context.Background(), uid); err != nil {
		handleResponse(c, "error while deleting staff tariff", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "staff tariff deleted")
}
