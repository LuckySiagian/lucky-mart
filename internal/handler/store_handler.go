package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService services.StoreService
}

func NewStoreHandler(storeService services.StoreService) *StoreHandler {
	return &StoreHandler{storeService: storeService}
}

// GetStore returns the store's public details (address, WhatsApp, hours).
func (h *StoreHandler) GetStore(c *gin.Context) {
	res, err := h.storeService.Get(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Info toko", res)
}

func (h *StoreHandler) UpdateStore(c *gin.Context) {
	var req dto.UpdateStoreSettingDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	res, err := h.storeService.Update(c.Request.Context(), req)
	if err != nil {
		sendServiceError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Info toko disimpan", res)
}
