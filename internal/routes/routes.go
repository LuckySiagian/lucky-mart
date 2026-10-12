package routes

import (
	"strings"

	"backend-sembako/config"
	_ "backend-sembako/docs"
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/internal/repositories"
	"backend-sembako/internal/services"
	customvalidator "backend-sembako/internal/validator"
	"backend-sembako/pkg/eventbus"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Register Custom DTO Validators
	customvalidator.RegisterCustomValidators()

	r := gin.New()
	r.MaxMultipartMemory = 8 << 20

	// Global Middleware Registration
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(cfg.CORS.AllowOrigins))
	r.Use(middleware.GlobalErrorHandler())

	// Initialize Shared EventBus
	eventBus := eventbus.NewEventBus()

	// Initialize Shared JWT Service
	jwtService := jwt.NewJWTService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(db)
	b2bProfileRepo := repositories.NewB2BProfileRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	loginHistoryRepo := repositories.NewLoginHistoryRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	stockMoveRepo := repositories.NewStockMovementRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	cartItemRepo := repositories.NewCartItemRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	invoiceRepo := repositories.NewInvoiceRepository(db)
	addressRepo := repositories.NewAddressRepository(db)
	courierRepo := repositories.NewCourierRepository(db)
	methodRepo := repositories.NewShippingMethodRepository(db)
	rateRepo := repositories.NewShippingRateRepository(db)
	shipmentRepo := repositories.NewShipmentRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	paymentMethodRepo := repositories.NewPaymentMethodRepository(db)
	paymentLogRepo := repositories.NewPaymentLogRepository(db)
	creditLedgerRepo := repositories.NewCreditLedgerRepository(db)
	creditPayRepo := repositories.NewCreditPaymentRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	templateRepo := repositories.NewNotificationTemplateRepository(db)
	notificationLogRepo := repositories.NewNotificationLogRepository(db)
	eventLogRepo := repositories.NewEventLogRepository(db)

	// Initialize Services
	healthService := services.NewHealthService(db)
	authService := services.NewAuthService(userRepo, b2bProfileRepo, refreshTokenRepo, loginHistoryRepo, jwtService, cfg.JWT.RefreshTTL, eventBus)
	userService := services.NewUserService(db, userRepo, b2bProfileRepo, refreshTokenRepo, eventBus)
	productService := services.NewProductService(db, productRepo, categoryRepo, inventoryRepo, stockMoveRepo)
	inventoryService := services.NewInventoryService(db, inventoryRepo, stockMoveRepo)
	cartService := services.NewCartService(db, cartRepo, cartItemRepo, productRepo, methodRepo, rateRepo)
	orderService := services.NewOrderService(db, orderRepo, invoiceRepo, cartRepo, b2bProfileRepo, addressRepo, paymentRepo, paymentMethodRepo, creditLedgerRepo, methodRepo, rateRepo, eventBus)
	addressService := services.NewAddressService(addressRepo)
	shippingService := services.NewShippingService(db, courierRepo, methodRepo, rateRepo, shipmentRepo, orderRepo, eventBus)
	paymentService := services.NewPaymentService(db, paymentRepo, paymentMethodRepo, invoiceRepo, orderRepo, paymentLogRepo, eventBus)
	paymentMethodService := services.NewPaymentMethodService(paymentMethodRepo)
	creditLedgerService := services.NewCreditLedgerService(db, creditLedgerRepo, creditPayRepo, userRepo, eventBus)
	notificationService := services.NewNotificationService(db, notificationRepo, notificationLogRepo)
	templateService := services.NewNotificationTemplateService(templateRepo)
	eventListenerService := services.NewEventListenerService(db, eventLogRepo, templateRepo, notificationRepo, notificationLogRepo)

	// Register EventBus Subscribers
	eventListenerService.RegisterSubscribers(eventBus)

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler(healthService, cfg)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	addressHandler := handler.NewAddressHandler(addressService)
	shippingHandler := handler.NewShippingHandler(shippingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	paymentMethodHandler := handler.NewPaymentMethodHandler(paymentMethodService)
	creditLedgerHandler := handler.NewCreditLedgerHandler(creditLedgerService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	templateHandler := handler.NewNotificationTemplateHandler(templateService)
	eventHandler := handler.NewEventHandlerController(eventListenerService)
	dashboardHandler := handler.NewDashboardHandler(services.NewDashboardService(db))
	uploadHandler := handler.NewUploadHandler(services.NewUploadService(cfg.Storage.UploadDir))
	storeHandler := handler.NewStoreHandler(services.NewStoreService(db))

	// Uploaded photos. Only the images folder is public; nosniff stops a
	// browser from treating a file as anything but the image it was saved as.
	images := r.Group(strings.TrimSuffix(services.ImagePathPrefix, "/"), func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	})
	images.Static("/", services.ImageDir(cfg.Storage.UploadDir))

	// Swagger Endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// System & Health Endpoints
	r.GET("/health/liveness", healthHandler.LivenessCheck)
	r.GET("/health/readiness", healthHandler.ReadinessCheck)
	r.GET("/version", healthHandler.VersionCheck)

	// API V1 Base Group
	v1 := r.Group("/api/v1")
	{
		RegisterAuthRoutes(v1, authHandler, jwtService)
		RegisterUserRoutes(v1, userHandler, jwtService)
		RegisterProductRoutes(v1, productHandler, jwtService)
		RegisterInventoryRoutes(v1, inventoryHandler, jwtService)
		RegisterCartRoutes(v1, cartHandler, jwtService)
		RegisterOrderRoutes(v1, orderHandler, jwtService)
		RegisterAddressRoutes(v1, addressHandler, jwtService)
		RegisterShippingRoutes(v1, shippingHandler, jwtService)
		RegisterPaymentRoutes(v1, paymentHandler, paymentMethodHandler, creditLedgerHandler, jwtService)
		RegisterNotificationRoutes(v1, notificationHandler, templateHandler, eventHandler, jwtService)

		dashboard := v1.Group("/admin/dashboard")
		dashboard.Use(middleware.JWTAuth(jwtService), middleware.RequireRoles(warehouseRole...))
		dashboard.GET("", dashboardHandler.GetStats)

		v1.GET("/store", storeHandler.GetStore)
		v1.PUT("/store", middleware.JWTAuth(jwtService), middleware.RequireRoles(managerRoles...), storeHandler.UpdateStore)

		uploads := v1.Group("/uploads")
		uploads.Use(middleware.JWTAuth(jwtService))
		uploads.POST("/images", uploadHandler.UploadImage)
	}

	return r
}
