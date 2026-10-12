package routes

import (
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
	methods := r.Group("/payment-methods")
	{
		methods.GET("", methodHandler.GetPaymentMethods)

		adminOnly := methods.Group("")
		adminOnly.Use(middleware.JWTAuth(jwtService), middleware.RequireRoles(managerRoles...))
		{
			adminOnly.GET("/all", methodHandler.GetAllPaymentMethods)
			adminOnly.POST("", methodHandler.CreatePaymentMethod)
			adminOnly.PUT("/:id", methodHandler.UpdatePaymentMethod)
			adminOnly.DELETE("/:id", methodHandler.DeletePaymentMethod)
		}
	}

	payments := r.Group("/payments")
	payments.Use(middleware.JWTAuth(jwtService))
	{
		// Customers: own payments only (enforced in the service).
		payments.GET("", paymentHandler.GetPayments)
		payments.POST("", paymentHandler.CreatePayment)
		payments.POST("/upload-proof", paymentHandler.UploadProof)
		payments.GET("/:id", paymentHandler.GetPaymentByID)

		// Verification is for managers only.
		finance := payments.Group("")
		finance.Use(middleware.RequireRoles(managerRoles...))
		{
			finance.GET("/waiting-confirmation", paymentHandler.GetWaitingConfirmation)
			finance.PATCH("/:id/approve", paymentHandler.ApprovePayment)
			finance.PATCH("/:id/reject", paymentHandler.RejectPayment)
			finance.PATCH("/:id/refund", paymentHandler.RefundPayment)
			finance.PATCH("/:id/cancel", paymentHandler.CancelPayment)
		}
	}

	ledgers := r.Group("/credit-ledgers")
	ledgers.Use(middleware.JWTAuth(jwtService))
	{
		// B2B customer: own kasbon account
		ledgers.GET("/me", ledgerHandler.GetMyLedger)
		ledgers.GET("/me/history", ledgerHandler.GetMyLedgerHistory)

		finance := ledgers.Group("")
		finance.Use(middleware.RequireRoles(managerRoles...))
		{
			finance.GET("", ledgerHandler.GetLedgers)
			finance.POST("", ledgerHandler.CreateLedger)
			finance.GET("/:id", ledgerHandler.GetLedgerByID)
			finance.PATCH("/:id", ledgerHandler.AdjustCreditLimit)
			finance.POST("/:id/settle", ledgerHandler.SettleCredit)
			finance.GET("/:id/history", ledgerHandler.GetPaymentHistory)
		}
	}
}
