package handler

import (
	"net/http"

	"backend-sembako/config"
	"backend-sembako/helpers"
	"backend-sembako/internal/services"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthService services.HealthService
	config        *config.Config
}

func NewHealthHandler(healthService services.HealthService, cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
		config:        cfg,
	}
}

// LivenessCheck godoc
// @Summary Liveness check endpoint
// @Description Returns 200 OK if server process is running
// @Tags System
// @Produce json
// @Success 200 {object} helpers.APIResponse
// @Router /health/liveness [get]
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	helpers.ResponseSuccess(c, http.StatusOK, "Liveness check passed", gin.H{
		"status": "alive",
	})
}

// ReadinessCheck godoc
// @Summary Readiness check endpoint
// @Description Returns 200 OK if database and downstream services are healthy
// @Tags System
// @Produce json
// @Success 200 {object} helpers.APIResponse
// @Failure 503 {object} helpers.APIErrorResponse
// @Router /health/readiness [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	isReady, msg := h.healthService.CheckReadiness(c.Request.Context())
	if !isReady {
		helpers.ResponseError(c, http.StatusServiceUnavailable, "Readiness check failed", gin.H{
			"status":  "not_ready",
			"details": msg,
		})
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Readiness check passed", gin.H{
		"status":   "ready",
		"database": msg,
	})
}

// VersionCheck godoc
// @Summary Application version endpoint
// @Description Returns application name, version, and current environment
// @Tags System
// @Produce json
// @Success 200 {object} helpers.APIResponse
// @Router /version [get]
func (h *HealthHandler) VersionCheck(c *gin.Context) {
	helpers.ResponseSuccess(c, http.StatusOK, "Version details", gin.H{
		"name":    h.config.App.Name,
		"version": h.config.App.Version,
		"env":     h.config.App.Env,
	})
}
