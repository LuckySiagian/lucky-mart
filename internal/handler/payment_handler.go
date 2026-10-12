package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(paymentService services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, role := getUserContext(c)
	var req dto.CreatePaymentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.paymentService.CreatePayment(c.Request.Context(), userID, role, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Payment record created successfully", res)
}

func (h *PaymentHandler) GetPayments(c *gin.Context) {
	userID, role := getUserContext(c)
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.paymentService.GetPayments(c.Request.Context(), userID, role, req, c.Query("status"))
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payments retrieved successfully", res)
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	userID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	res, err := h.paymentService.GetPaymentByID(c.Request.Context(), userID, role, id)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment details retrieved successfully", res)
}

func (h *PaymentHandler) UploadProof(c *gin.Context) {
	userID, _ := getUserContext(c)
	var req dto.UploadProofDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.paymentService.UploadProof(c.Request.Context(), userID, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Proof of transfer uploaded successfully", res)
}

func (h *PaymentHandler) GetWaitingConfirmation(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.paymentService.GetWaitingConfirmation(c.Request.Context(), req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payments waiting for confirmation retrieved successfully", res)
}

func (h *PaymentHandler) ApprovePayment(c *gin.Context) {
	adminID, _ := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ConfirmPaymentDTO
	_ = c.ShouldBindJSON(&req)

	res, err := h.paymentService.ApprovePayment(c.Request.Context(), adminID, id, req.Notes)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment approved successfully", res)
}

func (h *PaymentHandler) RejectPayment(c *gin.Context) {
	adminID, _ := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.RejectPaymentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.paymentService.RejectPayment(c.Request.Context(), adminID, id, req.Notes)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment rejected successfully", res)
}

func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	userID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	res, err := h.paymentService.CancelPayment(c.Request.Context(), userID, role, id)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment cancelled successfully", res)
}

func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	adminID, _ := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.RefundDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.paymentService.RefundPayment(c.Request.Context(), adminID, id, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment refunded successfully", res)
}
