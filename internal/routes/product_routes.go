package routes

import (
	"backend-sembako/internal/handler"
	"backend-sembako/internal/middleware"
	"backend-sembako/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.RouterGroup, productHandler *handler.ProductHandler, jwtService jwt.JWTService) {
	// Public catalog
	r.GET("/categories", productHandler.GetCategories)
	r.GET("/products", productHandler.GetProducts)
	r.GET("/products/:id", productHandler.GetProductByID)

	// Catalog management (managers only)
	manage := r.Group("")
	manage.Use(middleware.JWTAuth(jwtService), middleware.RequireRoles(managerRoles...))
	{
		manage.POST("/categories", productHandler.CreateCategory)
		manage.PUT("/categories/:id", productHandler.UpdateCategory)
		manage.DELETE("/categories/:id", productHandler.DeleteCategory)

		manage.POST("/products", productHandler.CreateProduct)
		manage.PUT("/products/:id", productHandler.UpdateProduct)
		manage.DELETE("/products/:id", productHandler.DeleteProduct)
	}

	// Staff view of the catalog including hidden products
	staff := r.Group("/admin")
	staff.Use(middleware.JWTAuth(jwtService), middleware.RequireRoles(warehouseRole...))
	{
		staff.GET("/products", productHandler.GetProductsAdmin)
	}
}
