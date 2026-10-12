package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, role := getUserContext(c)
	res, err := h.cartService.GetCart(c.Request.Context(), userID, role)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Cart retrieved successfully", res)
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, role := getUserContext(c)
	var req dto.AddToCartDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.cartService.AddToCart(c.Request.Context(), userID, role, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Item added to cart successfully", res)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, role := getUserContext(c)
	itemIDStr := c.Param("id")
	itemID, err := utils.ParseUUID(itemIDStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid cart item UUID format", nil)
		return
	}

	var req dto.UpdateCartItemDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.cartService.UpdateCartItem(c.Request.Context(), userID, role, itemID, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Cart item updated successfully", res)
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID, role := getUserContext(c)
	itemIDStr := c.Param("id")
	itemID, err := utils.ParseUUID(itemIDStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid cart item UUID format", nil)
		return
	}

	res, err := h.cartService.RemoveCartItem(c.Request.Context(), userID, role, itemID)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Cart item removed successfully", res)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, _ := getUserContext(c)
	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Cart cleared successfully", nil)
}

func (h *CartHandler) GetCheckoutSummary(c *gin.Context) {
	userID, role := getUserContext(c)
	var req dto.CheckoutSummaryRequestDTO
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(err)
			return
		}
	}
	res, err := h.cartService.GetCheckoutSummary(c.Request.Context(), userID, role, req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Checkout summary generated successfully", res)
}
