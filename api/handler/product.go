package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sell/api/models"
	"strconv"
)

// CreateProduct godoc
// @Router       /product [POST]
// @Summary      Create a new product
// @Description  create a new product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param 		 product body models.CreateProduct false "sale"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) CreateProduct(c *gin.Context) {
	product := models.CreateProduct{}

	if err := c.ShouldBindJSON(&product); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storage.Product().Create(context.Background(), product)
	if err != nil {
		handleResponse(c, "error is while creating product", http.StatusInternalServerError, err.Error())
		return
	}

	createdProduct, err := h.storage.Product().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusCreated, createdProduct)
}

// GetProduct godoc
// @Router       /product/{id} [GET]
// @Summary      Get product by id
// @Description  get product by id
// @Tags         product
// @Accept       json
// @Produce      json
// @Param 		 id path string true "product_id"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetProduct(c *gin.Context) {
	uid := c.Param("id")
	product, err := h.storage.Product().GetByID(context.Background(), uid)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, product)
}

// GetProductList godoc
// @Router       /products [GET]
// @Summary      Get product list
// @Description  get product list
// @Tags         product
// @Accept       json
// @Produce      json
// @Param 		 page query string false "page"
// @Param 		 limit query string false "limit"
// @Param 		 name query string false "name"
// @Param 		 barcode query int false "barcode"
// @Success      200  {object}  models.ProductResponse
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) GetProductList(c *gin.Context) {
	var (
		page, limit int
		name        string
		barcode     int
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

	name = c.Query("search")

	barcode, err = strconv.Atoi(c.DefaultQuery("barcode", "0"))
	if err != nil {
		handleResponse(c, "error is while converting barcode", http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.storage.Product().GetList(context.Background(), models.ProductGetListRequest{
		Page:    page,
		Limit:   limit,
		Name:    name,
		Barcode: barcode,
	})
	if err != nil {
		handleResponse(c, "error is while getting list", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, products)
}

// UpdateProduct godoc
// @Router       /product/{id} [PUT]
// @Summary      Update product
// @Description  update
// @Tags         product
// @Accept       json
// @Produce      json
// @Param 		 id path string true "product_id"
// @Param 		 product body models.UpdateProduct false "product"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) UpdateProduct(c *gin.Context) {
	uid := c.Param("id")
	product := models.UpdateProduct{}

	if err := c.ShouldBindJSON(&product); err != nil {
		handleResponse(c, "error is while reading body", http.StatusBadRequest, err.Error())
		return
	}

	product.ID = uid
	id, err := h.storage.Product().Update(context.Background(), product)
	if err != nil {
		handleResponse(c, "error is while updating", http.StatusInternalServerError, err.Error())
		return
	}

	updatedProduct, err := h.storage.Product().GetByID(context.Background(), id)
	if err != nil {
		handleResponse(c, "error is while getting by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, updatedProduct)
}

// DeleteProduct godoc
// @Router       /product/{id} [DELETE]
// @Summary      Delete product
// @Description  delete product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param 		 id path string true "product_id"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
func (h Handler) DeleteProduct(c *gin.Context) {
	uid := c.Param("id")
	if err := h.storage.Product().Delete(context.Background(), uid); err != nil {
		handleResponse(c, "error is while deleting", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "product deleted!")
}
