package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
	"strconv"
)

// CreateBranch godoc
// @Router       /branch [POST]
// @Summary      Create a new branch
// @Description  create a new branch
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param 		 branch body models.CreateBranch false "branch"
// @Success      200  {object}  models.Branch
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateBranch(c *gin.Context) {
	branch := models.CreateBranch{}

	if err := c.ShouldBindJSON(&branch); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Branch().Create(context.Background(), branch)
	if err != nil {
		handleResponse(c, "error is while creating branch", http.StatusInternalServerError, err.Error())
		return
	}

	createdBranch, err := h.storage.Branch().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdBranch)
}

// GetBranch godoc
// @Router       /branch/{id} [GET]
// @Summary      Get branch by id
// @Description  get branch by id
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param 		 id path string true "branch_id"
// @Success      200  {object}  models.Branch
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBranch(c *gin.Context) {
	uid := c.Param("id")

	branch, err := h.storage.Branch().GetByID(context.Background(), uid)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, branch)
}

// GetBranchList godoc
// @Router       /branches [GET]
// @Summary      Get branch list
// @Description  get branch list
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.BranchResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetBranchList(c *gin.Context) {
	var (
		page, limit int
		search      string
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error is while converting page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error is while converting limit", http.StatusBadRequest, err.Error())
		return
	}

	search = c.Query("search")

	branches, err := h.storage.Branch().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})

	handleResponse(c, "", http.StatusOK, branches)
}

// UpdateBranch godoc
// @Router       /branch/{id} [PUT]
// @Summary      Update branch
// @Description  update branch
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param 		 id path string true "branch_id"
// @Param 		 branch body models.UpdateBranch false "branch"
// @Success      200  {object}  models.Branch
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateBranch(c *gin.Context) {
	uid := c.Param("id")

	branch := models.UpdateBranch{}
	if err := c.ShouldBindJSON(&branch); err != nil {
		handleResponse(c, "error is wile reading from body", http.StatusBadRequest, err.Error())
		return
	}

	branch.ID = uid
	id, err := h.storage.Branch().Update(context.Background(), branch)
	if err != nil {
		handleResponse(c, "error is while updating branch", http.StatusInternalServerError, err.Error())
		return
	}

	updatedBranch, err := h.storage.Branch().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedBranch)
}

// DeleteBranch godoc
// @Router       /branch/{id} [DELETE]
// @Summary      Delete branch
// @Description  delete branch
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param 		 id path string true "branch_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteBranch(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.Branch().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error is while delteing branch", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(c, "", http.StatusOK, "branch deleted!")
}
