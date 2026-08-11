package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(
	router *gin.RouterGroup,
	productHandler *handlers.ProductHandler,
	jwtManager *utils.JWTManager,
) {
	products := router.Group("/products")
	{
		products.GET(
			"",
			productHandler.GetAllProducts,
		)

		products.GET(
			"/search",
			productHandler.SearchProducts,
		)

		products.GET(
			"/:id",
			productHandler.GetProductByID,
		)
	}

	adminProducts := router.Group("/admin/products")

	adminProducts.Use(
		middleware.AuthMiddleware(jwtManager),
		middleware.RequireRoles("owner"),
	)

	{
		adminProducts.POST(
			"",
			productHandler.CreateProduct,
		)

		adminProducts.PUT(
			"/:id",
			productHandler.UpdateProduct,
		)

		adminProducts.DELETE(
			"/:id",
			productHandler.DeleteProduct,
		)
	}
}