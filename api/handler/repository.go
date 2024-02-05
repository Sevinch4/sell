package handler

import (
	"context"
	"net/http"
	"sell/api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRepository godoc
// @Router       /repository [POST]
// @Summary      Create a new repository
// @Description  create a new repository
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param 		 repository body models.CreateRepository false "repository"
// @Success      200  {object}  models.Repository
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateRepository(c *gin.Context) {
	repository := models.CreateRepository{}

	if err := c.ShouldBindJSON(&repository); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Repository().Create(context.Background(), repository)
	if err != nil {
		handleResponse(c, "error while creating repository", http.StatusInternalServerError, err.Error())
		return
	}

	createdRepository, err := h.storage.Repository().GetByID(context.Background(), models.PrimaryKey{
		ID: id,
	})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdRepository)
}

// GetRepository godoc
// @Router       /repository/{id} [GET]
// @Summary      Get repository by id
// @Description  get repository by id
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param 		 id path string true "repository_id"
// @Success      200  {object}  models.Repository
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetRepository(c *gin.Context) {
	uid := c.Param("id")

	repository, err := h.storage.Repository().GetByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting repository by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, repository)
}

// GetRepositoryList godoc
// @Router       /repositories [GET]
// @Summary      Get repository list
// @Description  get repository list
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 search query string false "search"
// @Success      200  {object}  models.RepositoriesResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetRepositoryList(c *gin.Context) {
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

	response, err := h.storage.Repository().GetList(context.Background(), models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(c, "error while getting repository list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, response)
}

// UpdateRepository godoc
// @Router       /repository/{id} [PUT]
// @Summary      Update repository
// @Description  get repository
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param 		 id path string true "repository_id"
// @Param 		 repository body models.UpdateRepository false "repository"
// @Success      200  {object}  models.Repository
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateRepository(c *gin.Context) {
	uid := c.Param("id")

	repository := models.UpdateRepository{}
	if err := c.ShouldBindJSON(&repository); err != nil {
		handleResponse(c, "error while reading from body", http.StatusBadRequest, err.Error())
		return
	}

	repository.ID = uid
	if _, err := h.storage.Repository().Update(context.Background(), repository); err != nil {
		handleResponse(c, "error while updating repository ", http.StatusInternalServerError, err.Error())
		return
	}

	updatedRepository, err := h.storage.Repository().GetByID(context.Background(), models.PrimaryKey{ID: uid})
	if err != nil {
		handleResponse(c, "error while getting by ID", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedRepository)
}

// DeleteRepository godoc
// @Router       /repository/{id} [DELETE]
// @Summary      Delete repository
// @Description  delete repository
// @Tags         repository
// @Accept       json
// @Produce      json
// @Param 		 id path string true "repository_id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteRepository(c *gin.Context) {
	uid := c.Param("id")

	if err := h.storage.Repository().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error while deleting repository ", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "repository deleted")
}
