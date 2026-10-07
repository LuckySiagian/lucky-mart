package routes

import (
	"backend-sembako/internal/constants"
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(
	r *gin.RouterGroup,
	paymentHandler *handler.PaymentHandler,
	methodHandler *handler.PaymentMethodHandler,
	ledgerHandler *handler.CreditLedgerHandler,
	jwtService jwt.JWTService,
) {
	// Payment Methods Public & Admin
	methods := r.Group("/payment-methods")
	{
		methods.GET("", methodHandler.GetPaymentMethods)

		adminOnly := methods.Group("")
		adminOnly.Use(middleware.JWTAuth(jwtService))
		adminOnly.Use(middleware.RequireRoles(constants.RoleSuperAdmin, constants.RoleStoreManager))
		{
			adminOnly.POST("", methodHandler.CreatePaymentMethod)
			adminOnly.PUT("/:id", methodHandler.UpdatePaymentMethod)
			adminOnly.DELETE("/:id", methodHandler.DeletePaymentMethod)
		}
	}

	// Payments Protected
	payments := r.Group("/payments")
	payments.Use(middleware.JWTAuth(jwtService))
	{
		payments.GET("", paymentHandler.GetPayments)
		payments.POST("", paymentHandler.CreatePayment)
		payments.POST("/upload-proof", paymentHandler.UploadProof)
		payments.GET("/:id", paymentHandler.GetPaymentByID)
		payments.PATCH("/:id/confirm", paymentHandler.ApprovePayment)
		payments.PATCH("/:id/cancel", paymentHandler.CancelPayment)

		// Finance & Admin Verification Routes
		financeOnly := payments.Group("")
		financeOnly.Use(middleware.RequireRoles(constants.RoleSuperAdmin, constants.RoleStoreManager))
		{
			financeOnly.GET("/waiting-confirmation", paymentHandler.GetWaitingConfirmation)
			financeOnly.PATCH("/:id/approve", paymentHandler.ApprovePayment)
			financeOnly.PATCH("/:id/reject", paymentHandler.RejectPayment)
			financeOnly.PATCH("/:id/refund", paymentHandler.RefundPayment)
		}
	}

	// Credit Ledgers (B2B Admin & Finance)
	ledgers := r.Group("/credit-ledgers")
	ledgers.Use(middleware.JWTAuth(jwtService))
	ledgers.Use(middleware.RequireRoles(constants.RoleSuperAdmin, constants.RoleStoreManager))
	{
		ledgers.GET("", ledgerHandler.GetLedgers)
		ledgers.GET("/:id", ledgerHandler.GetLedgerByID)
		ledgers.POST("", ledgerHandler.CreateLedger)
		ledgers.PATCH("/:id", ledgerHandler.AdjustCreditLimit)
		ledgers.GET("/:id/history", ledgerHandler.GetPaymentHistory)
	}
}
