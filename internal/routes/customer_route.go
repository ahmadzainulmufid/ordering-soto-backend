package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterCustomerRoutes(
	router *gin.RouterGroup,
	customerHandler *handlers.CustomerAuthHandler,
	jwtManager *utils.JWTManager,
) {
	customer := router.Group("/customer/auth")
	{
		customer.POST("/register", customerHandler.Register)
		customer.POST("/login", customerHandler.Login)

		customer.GET(
			"/profile",
			middleware.AuthMiddleware(jwtManager),
			middleware.RequireRoles("customer"),
			customerHandler.Profile,
		)
	}

	// Contoh untuk endpoint khusus customer nanti:
	// orders := router.Group("/customer/orders")
	// orders.Use(
	// 	middleware.AuthMiddleware(jwtManager),
	// 	middleware.RequireRoles("customer"),
	// )
	// {
	// 	orders.POST("", orderHandler.Create)
	// 	orders.GET("", orderHandler.ListMine)
	// }
}