package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService services.InventoryService
}

func NewInventoryHandler(inventoryService services.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	var req dto.StockAdjustmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.inventoryService.AdjustStock(c.Request.Context(), req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
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

func (h *InventoryHandler) GetStockHistory(c *gin.Context) {
	inventoryIDStr := c.Param("id")
	inventoryID, err := utils.ParseUUID(inventoryIDStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid inventory UUID format", nil)
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
