package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService services.InventoryService
}

func NewInventoryHandler(inventoryService services.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) GetInventories(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.inventoryService.GetInventories(c.Request.Context(), req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Inventory retrieved successfully", res)
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	var req dto.StockAdjustmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.inventoryService.AdjustStock(c.Request.Context(), req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Stock adjustment processed successfully", res)
}

func (h *InventoryHandler) GetLowStockAlerts(c *gin.Context) {
	res, err := h.inventoryService.GetLowStockAlerts(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Low stock items retrieved successfully", res)
}

func (h *InventoryHandler) GetMovements(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.inventoryService.GetMovements(c.Request.Context(), req, c.Query("type"))
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Stock movements retrieved successfully", res)
}

func (h *InventoryHandler) GetStockHistory(c *gin.Context) {
	inventoryID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.inventoryService.GetStockHistory(c.Request.Context(), inventoryID, req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Stock movement history retrieved successfully", res)
}

func (h *InventoryHandler) UpdateMinStockAlert(c *gin.Context) {
	inventoryID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateMinStockDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.inventoryService.UpdateMinStockAlert(c.Request.Context(), inventoryID, req.MinStockAlert)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Minimum stock alert updated successfully", res)
}
