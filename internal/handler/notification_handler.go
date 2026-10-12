package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notifService services.NotificationService
}

func NewNotificationHandler(notifService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: notifService}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, role := getUserRoleInfo(c)
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.notifService.GetNotifications(c.Request.Context(), userID, role, req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notifications retrieved successfully", res)
}

func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	userID, role := getUserRoleInfo(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid notification UUID format", nil)
		return
	}

	res, err := h.notifService.GetNotificationByID(c.Request.Context(), userID, role, id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notification details retrieved successfully", res)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, _ := getUserRoleInfo(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid notification UUID format", nil)
		return
	}

	res, err := h.notifService.MarkAsRead(c.Request.Context(), userID, id)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notification marked as read", res)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, _ := getUserRoleInfo(c)
	if err := h.notifService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "All notifications marked as read", nil)
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, _ := getUserRoleInfo(c)
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid notification UUID format", nil)
		return
	}

	if err := h.notifService.DeleteNotification(c.Request.Context(), userID, id); err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notification deleted successfully", nil)
}

func (h *NotificationHandler) RetryNotification(c *gin.Context) {
	var req dto.RetryNotificationDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.notifService.RetryNotification(c.Request.Context(), req.NotificationID)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notification retry attempt executed successfully", res)
}

func getUserRoleInfo(c *gin.Context) (uuid.UUID, string) {
	u, roleStr := getUserContextInfo(c)
	return u, roleStr
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := getUserContext(c)
	count, err := h.notifService.CountUnread(c.Request.Context(), userID)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Unread notification count retrieved", gin.H{"unread": count})
}
