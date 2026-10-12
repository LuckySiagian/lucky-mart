package handler

import (
	"net/http"

	"backend-sembako/helpers"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(dashboardService services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	res, err := h.dashboardService.GetStats(c.Request.Context())
	if err != nil {
		helpers.SendInternalServerError(c, err.Error())
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Dashboard statistics retrieved successfully", res)
}
