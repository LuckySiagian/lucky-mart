package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// GetProducts lists the public catalog. With include_inactive=true (admin
// route only) hidden products are listed too.
func (h *ProductHandler) GetProducts(c *gin.Context) {
	h.listProducts(c, false)
}

func (h *ProductHandler) GetProductsAdmin(c *gin.Context) {
	h.listProducts(c, true)
}

func (h *ProductHandler) listProducts(c *gin.Context, includeInactive bool) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	var categoryID *uuid.UUID
	if catStr := c.Query("category_id"); catStr != "" {
		if id, err := utils.ParseUUID(catStr); err == nil {
			categoryID = &id
		}
	}

	res, err := h.productService.GetProducts(c.Request.Context(), req, categoryID, includeInactive)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Product catalog fetched successfully", res)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	res, err := h.productService.GetProductByID(c.Request.Context(), id, false)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Product details retrieved successfully", res)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productService.CreateProduct(c.Request.Context(), req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Product created successfully", res)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productService.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Product updated successfully", res)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetCategories(c *gin.Context) {
	res, err := h.productService.GetCategories(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Categories fetched successfully", res)
}

func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productService.CreateCategory(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Category created successfully", res)
}

func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.CreateCategoryDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productService.UpdateCategory(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Category updated successfully", res)
}

func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.productService.DeleteCategory(c.Request.Context(), id); err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Category deleted successfully", nil)
}
