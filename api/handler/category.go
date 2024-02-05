package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
	"strconv"
)

// CreateCategory godoc
// @Router       /category [POST]
// @Summary      Create a new category
// @Description  create a new category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param 		 category body models.CreateCategory false "category"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateCategory(c *gin.Context) {
	category := models.CreateCategory{}
	if err := c.ShouldBindJSON(&category); err != nil {
		handleResponse(c, "error is while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Category().Create(c.Request.Context(), category)
	if err != nil {
		handleResponse(c, "error is while creating category", http.StatusInternalServerError, err.Error())
		return
	}

	createdCategory, err := h.storage.Category().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdCategory)
}

// GetCategory godoc
// @Router       /category/{id} [GET]
// @Summary      Get category by id
// @Description  get category by id
// @Tags         category
// @Accept       json
// @Produce      json
// @Param 		 id path string true "category_id"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetCategory(c *gin.Context) {
	uid := c.Param("id")
	category, err := h.storage.Category().GetByID(context.Background(), uid)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, category)
}

// GetCategoryList godoc
// @Router       /categories [GET]
// @Summary      Get category list
// @Description  get category list
// @Tags         category
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.CategoryResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetCategoryList(c *gin.Context) {
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

	categories, err := h.storage.Category().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, "error is while getting list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, categories)
}

// UpdateCategory godoc
// @Router       /category/{id} [PUT]
// @Summary      Update category
// @Description  get category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param 		 id path string true "category_id"
// @Param 		 category body models.UpdateCategory false "category"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateCategory(c *gin.Context) {
	uid := c.Param("id")
	category := models.UpdateCategory{}

	if err := c.ShouldBindJSON(&category); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	category.ID = uid

	id, err := h.storage.Category().Update(context.Background(), category)
	if err != nil {
		handleResponse(c, "error is while updating category", http.StatusInternalServerError, err.Error())
		return
	}

	updatedCategory, err := h.storage.Category().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedCategory)
}

// DeleteCategory godoc
// @Router       /category/{id} [DELETE]
// @Summary      Delete category
// @Description  delete category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param 		 id path string true "category_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteCategory(c *gin.Context) {
	uid := c.Param("id")
	if err := h.storage.Category().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error is while deleting", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "category deleted!")
}
