package handler

import (
	"errors"

	"backend-sembako/helpers"
	"backend-sembako/internal/constants"
	"backend-sembako/internal/middleware"
	"backend-sembako/internal/services"
	"backend-sembako/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getUserContext returns the authenticated user's id and role, set by the
// JWTAuth middleware.
func getUserContext(c *gin.Context) (uuid.UUID, constants.UserRole) {
	userIDVal, _ := c.Get(middleware.ContextUserIDKey)
	roleVal, _ := c.Get(middleware.ContextUserRole)

	userID, _ := userIDVal.(uuid.UUID)
	role, _ := roleVal.(constants.UserRole)
	return userID, role
}

// getUserContextInfo is getUserContext with the role as a plain string.
func getUserContextInfo(c *gin.Context) (uuid.UUID, string) {
	userID, role := getUserContext(c)
	return userID, string(role)
}

// parseIDParam reads a UUID path parameter, answering 400 when malformed.
func parseIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := utils.ParseUUID(c.Param(name))
	if err != nil {
		helpers.SendBadRequest(c, "ID tidak valid", nil)
		return uuid.Nil, false
	}
	return id, true
}

var notFoundErrors = []error{
	services.ErrOrderNotFound, services.ErrPaymentNotFound, services.ErrProductNotFound,
	services.ErrCategoryNotFound, services.ErrUserNotFound, services.ErrShipmentNotFound,
	services.ErrLedgerNotFound, services.ErrInventoryNotFound, services.ErrAddressNotFound,
	services.ErrProductUnitNotFound, services.ErrPaymentMethodNotFound,
}

var forbiddenErrors = []error{
	services.ErrForbiddenUserEdit, services.ErrUnauthorizedPayment, services.ErrUnauthorizedOrder,
	services.ErrUnauthorizedShipment, services.ErrUnauthorizedAddress,
}

// sendServiceError maps service errors to HTTP status codes: not-found to
// 404, permission errors to 403, everything else (business-rule violations)
// to 400 with the service's message.
func sendServiceError(c *gin.Context, err error) {
	for _, target := range notFoundErrors {
		if errors.Is(err, target) {
			helpers.SendNotFound(c, err.Error())
			return
		}
	}
	for _, target := range forbiddenErrors {
		if errors.Is(err, target) {
			helpers.SendForbidden(c, err.Error())
			return
		}
	}
	helpers.SendBadRequest(c, err.Error(), nil)
}
