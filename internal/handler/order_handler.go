package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"

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
		sendServiceError(c, err)
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

	res, err := h.orderService.GetOrders(c.Request.Context(), userID, role, req, c.Query("status"))
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Orders retrieved successfully", res)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	res, err := h.orderService.GetOrderByID(c.Request.Context(), userID, role, id)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order details retrieved successfully", res)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.CancelOrderDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.orderService.CancelOrder(c.Request.Context(), userID, role, id, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order cancelled successfully and stock restored", res)
}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	userID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	res, err := h.orderService.GetOrderHistory(c.Request.Context(), userID, role, id)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order history retrieved successfully", res)
}

func (h *OrderHandler) GetInvoiceByID(c *gin.Context) {
	userID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
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
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateOrderStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.orderService.UpdateOrderStatus(c.Request.Context(), userID, role, id, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Order status updated successfully", res)
}
