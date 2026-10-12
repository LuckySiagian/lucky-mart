package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(
	r *gin.RouterGroup,
	notifHandler *handler.NotificationHandler,
	tmplHandler *handler.NotificationTemplateHandler,
	eventHandler *handler.EventHandlerController,
	jwtService jwt.JWTService,
) {
	// Notifications (User Protected)
	notifications := r.Group("/notifications")
	notifications.Use(middleware.JWTAuth(jwtService))
	{
		notifications.GET("", notifHandler.GetNotifications)
		notifications.GET("/unread-count", notifHandler.GetUnreadCount)
		notifications.PATCH("/read-all", notifHandler.MarkAllAsRead)
		notifications.GET("/:id", notifHandler.GetNotificationByID)
		notifications.PATCH("/:id/read", notifHandler.MarkAsRead)
		notifications.DELETE("/:id", notifHandler.DeleteNotification)
	}

	// Templates (Super Admin Protected)
	templates := r.Group("/notification-templates")
	templates.Use(middleware.JWTAuth(jwtService))
	templates.Use(middleware.RequireRoles(managerRoles...))
	{
		templates.GET("", tmplHandler.GetTemplates)
		templates.POST("", tmplHandler.CreateTemplate)
		templates.PUT("/:id", tmplHandler.UpdateTemplate)
		templates.DELETE("/:id", tmplHandler.DeleteTemplate)
	}

	// Domain Events (Admin Protected)
	events := r.Group("/events")
	events.Use(middleware.JWTAuth(jwtService))
	events.Use(middleware.RequireRoles(managerRoles...))
	{
		events.GET("", eventHandler.GetEvents)
		events.GET("/:id", eventHandler.GetEventByID)
		events.POST("/retry", notifHandler.RetryNotification)
	}
}
