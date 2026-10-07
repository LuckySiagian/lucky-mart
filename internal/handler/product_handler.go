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

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	var categoryID *uuid.UUID
	catStr := c.Query("category_id")
	if catStr != "" {
		if id, err := utils.ParseUUID(catStr); err == nil {
			categoryID = &id
		}
	}

	res, err := h.productService.GetProducts(c.Request.Context(), req, categoryID)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Product catalog fetched successfully", res)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid product UUID format", nil)
		return
	}

	res, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
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
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Product created successfully", res)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")

	var req dto.UpdateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.productService.UpdateProduct(c.Request.Context(), idStr, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Product updated successfully", res)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid product UUID format", nil)
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Product soft-deleted successfully", nil)
}

func (h *ProductHandler) GetCategories(c *gin.Context) {
	res, err := h.productService.GetCategories(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Categories fetched successfully", res)
}

type createCategoryReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.SendBadRequest(c, "Name is required", nil)
		return
	}

	res, err := h.productService.CreateCategory(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Category created successfully", res)
}

func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.SendBadRequest(c, "Name is required", nil)
		return
	}

	res, err := h.productService.UpdateCategory(c.Request.Context(), idStr, req.Name, req.Description)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Category updated successfully", res)
}

func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid category UUID format", nil)
		return
	}

	if err := h.productService.DeleteCategory(c.Request.Context(), id); err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Category deleted successfully", nil)
}
