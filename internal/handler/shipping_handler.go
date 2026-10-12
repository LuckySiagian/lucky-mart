package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type ShippingHandler struct {
	shippingService services.ShippingService
}

func NewShippingHandler(shippingService services.ShippingService) *ShippingHandler {
	return &ShippingHandler{shippingService: shippingService}
}

func (h *ShippingHandler) GetShippingMethods(c *gin.Context) {
	methods, err := h.shippingService.GetShippingMethods(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipping methods retrieved successfully", methods)
}

func (h *ShippingHandler) CalculateShippingRate(c *gin.Context) {
	var req dto.CalculateShippingDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.shippingService.CalculateShippingRate(c.Request.Context(), req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipping rate calculated successfully", res)
}

func (h *ShippingHandler) GetShipments(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.shippingService.GetShipments(c.Request.Context(), req, c.Query("status"))
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipments retrieved successfully", res)
}

func (h *ShippingHandler) CreateShipment(c *gin.Context) {
	var req dto.CreateShipmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.shippingService.CreateShipment(c.Request.Context(), req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Shipment created successfully", res)
}

func (h *ShippingHandler) GetShipmentByOrderID(c *gin.Context) {
	userID, role := getUserContext(c)
	orderID, ok := parseIDParam(c, "order_id")
	if !ok {
		return
	}

	res, err := h.shippingService.GetShipmentByOrderID(c.Request.Context(), userID, role, orderID)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipment tracking details retrieved successfully", res)
}

func (h *ShippingHandler) UpdateShipmentStatus(c *gin.Context) {
	shipmentID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateShipmentStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.shippingService.UpdateShipmentStatus(c.Request.Context(), shipmentID, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipment status updated successfully", res)
}
