package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

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
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipping rate calculated successfully", res)
}

func (h *ShippingHandler) CreateShipment(c *gin.Context) {
	var req dto.CreateShipmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.shippingService.CreateShipment(c.Request.Context(), req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Shipment created successfully", res)
}

func (h *ShippingHandler) GetShipmentByOrderID(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := utils.ParseUUID(orderIDStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid order UUID format", nil)
		return
	}

	res, err := h.shippingService.GetShipmentByOrderID(c.Request.Context(), orderID)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipment details retrieved successfully", res)
}

func (h *ShippingHandler) UpdateShipmentStatus(c *gin.Context) {
	shipmentIDStr := c.Param("id")
	shipmentID, err := utils.ParseUUID(shipmentIDStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid shipment UUID format", nil)
		return
	}

	var req dto.UpdateShipmentStatusDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.shippingService.UpdateShipmentStatus(c.Request.Context(), shipmentID, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Shipment status updated successfully", res)
}
