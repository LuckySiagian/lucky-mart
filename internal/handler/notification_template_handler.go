package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
)

type NotificationTemplateHandler struct {
	tmplService services.NotificationTemplateService
}

func NewNotificationTemplateHandler(tmplService services.NotificationTemplateService) *NotificationTemplateHandler {
	return &NotificationTemplateHandler{tmplService: tmplService}
}

func (h *NotificationTemplateHandler) GetTemplates(c *gin.Context) {
	res, err := h.tmplService.GetTemplates(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notification templates retrieved successfully", res)
}

func (h *NotificationTemplateHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateNotificationTemplateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.tmplService.CreateTemplate(c.Request.Context(), req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Notification template created successfully", res)
}

func (h *NotificationTemplateHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid template UUID format", nil)
		return
	}

	var req dto.CreateNotificationTemplateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.tmplService.UpdateTemplate(c.Request.Context(), id, req)
	if err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notification template updated successfully", res)
}

func (h *NotificationTemplateHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid template UUID format", nil)
		return
	}

	if err := h.tmplService.DeleteTemplate(c.Request.Context(), id); err != nil {
		helpers.SendBadRequest(c, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Notification template deleted successfully", nil)
}
