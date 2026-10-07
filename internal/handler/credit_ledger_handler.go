package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
)

type CreditLedgerHandler struct {
	ledgerService services.CreditLedgerService
}

func NewCreditLedgerHandler(ledgerService services.CreditLedgerService) *CreditLedgerHandler {
	return &CreditLedgerHandler{ledgerService: ledgerService}
}

func (h *CreditLedgerHandler) GetLedgers(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.ledgerService.GetLedgers(c.Request.Context(), req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Credit ledgers retrieved successfully", res)
}

func (h *CreditLedgerHandler) GetLedgerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid ledger UUID format", nil)
		return
	}

	res, err := h.ledgerService.GetLedgerByID(c.Request.Context(), id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Credit ledger details retrieved successfully", res)
}

func (h *CreditLedgerHandler) CreateLedger(c *gin.Context) {
	var req dto.CreateCreditLedgerDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.ledgerService.CreateLedger(c.Request.Context(), req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Credit ledger created successfully", res)
}

func (h *CreditLedgerHandler) AdjustCreditLimit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid ledger UUID format", nil)
		return
	}

	var req dto.AdjustCreditLimitDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.ledgerService.AdjustCreditLimit(c.Request.Context(), id, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Credit limit adjusted successfully", res)
}

func (h *CreditLedgerHandler) GetPaymentHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid ledger UUID format", nil)
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.ledgerService.GetPaymentHistory(c.Request.Context(), id, req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Credit payment history retrieved successfully", res)
}
