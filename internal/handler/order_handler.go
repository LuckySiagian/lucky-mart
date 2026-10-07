package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, role := getUserContext(c)
	var req dto.CreateOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.orderService.CreateOrder(c.Request.Context(), userID, role, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Order created successfully", res)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, role := getUserContext(c)
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.orderService.GetOrders(c.Request.Context(), userID, role, req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Orders retrieved successfully", res)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID, role := getUserContext(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid order UUID format", nil)
		return
	}

	res, err := h.orderService.GetOrderByID(c.Request.Context(), userID, role, id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order details retrieved successfully", res)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, role := getUserContext(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid order UUID format", nil)
		return
	}

	var req dto.CancelOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.orderService.CancelOrder(c.Request.Context(), userID, role, id, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order cancelled successfully and stock restored", res)
}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	userID, role := getUserContext(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid order UUID format", nil)
		return
	}

	res, err := h.orderService.GetOrderHistory(c.Request.Context(), userID, role, id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order history retrieved successfully", res)
}

func (h *OrderHandler) GetInvoiceByID(c *gin.Context) {
	userID, role := getUserContext(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid invoice UUID format", nil)
		return
	}

	res, err := h.orderService.GetInvoiceByID(c.Request.Context(), userID, role, id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Invoice retrieved successfully", res)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	userID, role := getUserContext(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid order UUID format", nil)
		return
	}

	var req dto.UpdateOrderStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.orderService.UpdateOrderStatus(c.Request.Context(), userID, role, id, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order status updated successfully", res)
}

