package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.RouterGroup, productHandler *handler.ProductHandler, jwtService jwt.JWTService) {
	r.GET("/categories", productHandler.GetCategories)
	r.POST("/categories", productHandler.CreateCategory)
	r.PUT("/categories/:id", productHandler.UpdateCategory)
	r.DELETE("/categories/:id", productHandler.DeleteCategory)

	products := r.Group("/products")
	{
		// Public & Admin Catalog Endpoints
		products.GET("", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.POST("", productHandler.CreateProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}
}
