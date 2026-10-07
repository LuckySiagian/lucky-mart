package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/dto"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
)

type EventHandlerController struct {
	listenerService services.EventListenerService
}

func NewEventHandlerController(listenerService services.EventListenerService) *EventHandlerController {
	return &EventHandlerController{listenerService: listenerService}
}

func (h *EventHandlerController) GetEvents(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.listenerService.GetEvents(c.Request.Context(), req)
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Domain Event logs retrieved successfully", res)
}

func (h *EventHandlerController) GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUUID(idStr)
	if err != nil {
		helpers.SendBadRequest(c, "Invalid event UUID format", nil)
		return
	}

	res, err := h.listenerService.GetEventByID(c.Request.Context(), id)
	if err != nil {
		helpers.SendNotFound(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Domain Event details retrieved successfully", res)
}
