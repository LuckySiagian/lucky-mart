package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/middleware"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressHandler struct {
	addressService services.AddressService
}

func NewAddressHandler(addressService services.AddressService) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID, _ := getUserContextInfo(c)
	var req dto.CreateAddressDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.addressService.CreateAddress(c.Request.Context(), userID, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Address created successfully", res)
}

func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userID, _ := getUserContextInfo(c)
	res, err := h.addressService.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Addresses retrieved successfully", res)
}

func (h *AddressHandler) GetAddressByID(c *gin.Context) {
	userID, _ := getUserContextInfo(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid address UUID format", nil)
		return
	}

	res, err := h.addressService.GetAddressByID(c.Request.Context(), userID, id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Address details retrieved successfully", res)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID, _ := getUserContextInfo(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid address UUID format", nil)
		return
	}

	var req dto.CreateAddressDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.addressService.UpdateAddress(c.Request.Context(), userID, id, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Address updated successfully", res)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID, _ := getUserContextInfo(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid address UUID format", nil)
		return
	}

	if err := h.addressService.DeleteAddress(c.Request.Context(), userID, id); err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Address deleted successfully", nil)
}

func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID, _ := getUserContextInfo(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid address UUID format", nil)
		return
	}

	res, err := h.addressService.SetDefaultAddress(c.Request.Context(), userID, id)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Primary default address updated successfully", res)
}

func getUserContextInfo(c *gin.Context) (uuid.UUID, string) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDVal.(uuid.UUID)
	return userID, ""
}
