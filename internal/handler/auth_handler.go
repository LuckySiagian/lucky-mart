package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/middleware"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterCustomerDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.authService.RegisterCustomer(c.Request.Context(), req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "User registered successfully", res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	device := c.GetHeader("X-Device-Info")
	if device == "" {
		device = "Unknown Device"
	}

	res, err := h.authService.Login(c.Request.Context(), req, ipAddress, userAgent, device)
	if err != nil {
		helpers.SendUnauthorized(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Authentication successful", res)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		helpers.SendUnauthorized(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Tokens refreshed successfully", res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		helpers.SendUnauthorized(c, "User not authenticated")
		return
	}

	userID := userIDVal.(uuid.UUID)
	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		helpers.SendInternalServerError(c, "Failed to revoke tokens")
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Logged out successfully", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		helpers.SendUnauthorized(c, "User not authenticated")
		return
	}

	userID := userIDVal.(uuid.UUID)
	profile, err := h.authService.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User profile retrieved successfully", profile)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDVal, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		helpers.SendUnauthorized(c, "User not authenticated")
		return
	}

	var req dto.ChangePasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	userID := userIDVal.(uuid.UUID)
	if err := h.authService.ChangePassword(c.Request.Context(), userID, req); err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Password changed successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), req); err != nil {
		helpers.SendInternalServerError(c, "Gagal memproses permintaan")
		return
	}

	// Same answer whether or not the number is registered.
	helpers.ResponseSuccess(c, http.StatusOK,
		"Permintaan reset kata sandi diterima. Hubungi admin toko untuk verifikasi identitas dan mendapatkan kata sandi baru.", nil)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := getUserContext(c)
	var req dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.authService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Profile updated successfully", res)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req); err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Password reset successfully. Please login with your new password.", nil)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), req); err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Email verified successfully", nil)
}
