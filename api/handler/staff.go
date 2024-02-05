package handler

import (
	"context"
	"net/http"
	"sell/api/models"
	"sell/pkg/check"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateStaff godoc
// @Router       /staff [POST]
// @Summary      Create a new staff
// @Description  create a new staff
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param 		 staff body models.CreateStaff false "staff"
// @Success      200  {object}  models.Staff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateStaff(c *gin.Context) {
	staff := models.CreateStaff{}

	if err := c.ShouldBindJSON(&staff); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Staff().Create(context.Background(), staff)

	if err != nil {
		handleResponse(c, "error while creating staff ", http.StatusInternalServerError, err.Error())
		return
	}

	createdStaffTarif, err := h.storage.Staff().StaffByID(context.Background(), models.PrimaryKey{ID: id})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdStaffTarif)
}

// GetStaff godoc
// @Router       /staff/{id} [GET]
// @Summary      Get staff by id
// @Description  get staff by id
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param 		 id path string true "staff_id"
// @Success      200  {object}  models.Staff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStaff(c *gin.Context) {
	uid := c.Param("id")

	staffTarif, err := h.storage.Staff().StaffByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting staff  by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, staffTarif)
}

// GetStaffList godoc
// @Router       /staffs [GET]
// @Summary      Get staff list
// @Description  get staff list
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.StaffsResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetStaffList(c *gin.Context) {
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

	response, err := h.storage.Staff().GetStaffTList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, "error while getting staff list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)
}

// UpdateStaff godoc
// @Router       /staff/{id} [PUT]
// @Summary      Update staff
// @Description  get staff
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param 		 id path string true "staff_id"
// @Param 		 staff body models.UpdateStaff false "staff"
// @Success      200  {object}  models.Staff
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateStaff(c *gin.Context) {
	uid := c.Param("id")

	staff := models.UpdateStaff{}
	if err := c.ShouldBindJSON(&staff); err != nil {
		handleResponse(c, "error while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	staff.ID = uid
	if _, err := h.storage.Staff().UpdateStaff(context.Background(), staff); err != nil {
		handleResponse(c, "error while updating staff ", http.StatusInternalServerError, err.Error())
		return
	}

	updatedStaff, err := h.storage.Staff().StaffByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedStaff)
}

// DeleteStaff godoc
// @Router       /staff/{id} [DELETE]
// @Summary      Delete staff
// @Description  delete staff
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param 		 id path string true "staff_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteStaff(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.Staff().DeleteStaff(context.Background(), uid); err != nil {
		handleResponse(c, "error while deleting staff ", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "staff deleted")
}

// UpdateStaffPassword godoc
// @Router       /staff/{id} [PATCH]
// @Summary      Update staff password
// @Description  update staff password
// @Tags         staff
// @Accept       json
// @Produce      json
// @Param 		 id path string true "staff_id"
// @Param        staff body models.UpdateStaffPassword true "staff"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateStaffPassword(c *gin.Context) {
	updateStaffPassword := models.UpdateStaffPassword{}

	if err := c.ShouldBindJSON(&updateStaffPassword); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		handleResponse(c, "error while parsing uuid", http.StatusBadRequest, err.Error())
		return
	}

	updateStaffPassword.ID = uid.String()

	oldPassword, err := h.storage.Staff().GetPassword(context.Background(), updateStaffPassword.ID)
	if err != nil {
		handleResponse(c, "error while getting password by id", http.StatusInternalServerError, err.Error())
		return
	}

	if oldPassword != updateStaffPassword.OldPassword {
		handleResponse(c, "old password is not correct", http.StatusBadRequest, "old password is not correct")
		return
	}

	if err = check.ValidatePassword(updateStaffPassword.NewPassword); err != nil {
		handleResponse(c, "new password is weak", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.Staff().UpdatePassword(context.Background(), updateStaffPassword); err != nil {
		handleResponse(c, "error while updating staff password by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "password successfully updated")
}
