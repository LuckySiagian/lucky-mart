package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
)

type PaymentMethodHandler struct {
	methodService services.PaymentMethodService
}

func NewPaymentMethodHandler(methodService services.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{methodService: methodService}
}

func (h *PaymentMethodHandler) GetPaymentMethods(c *gin.Context) {
	res, err := h.methodService.GetPaymentMethods(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment methods retrieved successfully", res)
}

func (h *PaymentMethodHandler) CreatePaymentMethod(c *gin.Context) {
	var req dto.CreatePaymentMethodDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.methodService.CreatePaymentMethod(c.Request.Context(), req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Payment method created successfully", res)
}

func (h *PaymentMethodHandler) UpdatePaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid payment method UUID format", nil)
		return
	}

	var req dto.CreatePaymentMethodDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.methodService.UpdatePaymentMethod(c.Request.Context(), id, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment method updated successfully", res)
}

func (h *PaymentMethodHandler) DeletePaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid payment method UUID format", nil)
		return
	}

	if err := h.methodService.DeletePaymentMethod(c.Request.Context(), id); err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment method deleted successfully", nil)
}
