package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentMethodHandler struct {
	methodService services.PaymentMethodService
}

func NewPaymentMethodHandler(methodService services.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{methodService: methodService}
}

// GetPaymentMethods lists the active methods customers can choose.
func (h *PaymentMethodHandler) GetPaymentMethods(c *gin.Context) {
	h.list(c, false)
}

// GetAllPaymentMethods (managers) also lists deactivated methods.
func (h *PaymentMethodHandler) GetAllPaymentMethods(c *gin.Context) {
	h.list(c, true)
}

func (h *PaymentMethodHandler) list(c *gin.Context, includeInactive bool) {
	res, err := h.methodService.GetPaymentMethods(c.Request.Context(), includeInactive)
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
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Payment method created successfully", res)
}

func (h *PaymentMethodHandler) UpdatePaymentMethod(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdatePaymentMethodDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.methodService.UpdatePaymentMethod(c.Request.Context(), id, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment method updated successfully", res)
}

func (h *PaymentMethodHandler) DeletePaymentMethod(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.methodService.DeletePaymentMethod(c.Request.Context(), id); err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment method deleted successfully", nil)
}
