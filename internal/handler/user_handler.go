package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.userService.GetUsers(c.Request.Context(), req, c.Query("role"))
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User list retrieved successfully", res)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	res, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User retrieved successfully", res)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	_, role := getUserContext(c)
	var req dto.CreateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.userService.CreateUser(c.Request.Context(), role, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "User created successfully", res)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	actorID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.userService.UpdateUser(c.Request.Context(), actorID, role, id, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User updated successfully", res)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	_, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.AdminResetPasswordDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), role, id, req.NewPassword); err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Password reset successfully", nil)
}

func (h *UserHandler) VerifyB2B(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.VerifyB2BDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.userService.VerifyB2B(c.Request.Context(), id, req)
	if err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "B2B verification updated successfully", res)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	actorID, role := getUserContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), actorID, role, id); err != nil {
		sendServiceError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User deleted successfully", nil)
}
