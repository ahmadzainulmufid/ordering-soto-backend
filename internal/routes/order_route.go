package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(
	router *gin.RouterGroup,
	orderHandler *handlers.OrderHandler,
	jwtManager *utils.JWTManager,
) {
	orders := router.Group("/orders")
	{
		orders.POST(
			"",
			orderHandler.CreateOrder,
		)

		orders.GET(
			"/code/:code",
			orderHandler.GetOrderByCode,
		)
	}

	adminOrders := router.Group("/admin/orders")

	adminOrders.Use(
		middleware.AuthMiddleware(jwtManager),
		middleware.RequireRoles(
			"owner",
			"cashier",
			"admin",
			"kitchen",
		),
	)

	{
		adminOrders.GET(
			"",
			orderHandler.GetAllOrders,
		)

		adminOrders.GET(
			"/:id",
			orderHandler.GetOrderByID,
		)

		adminOrders.PATCH(
			"/:id/status",
			orderHandler.UpdateOrderStatus,
		)
	}
}