package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterDiningTableRoutes(
	api *gin.RouterGroup,
	diningTableHandler *handlers.DiningTableHandler,
	jwtManager *utils.JWTManager,
) {
	// Public:
	// digunakan ketika customer scan QR meja.
	tables := api.Group("/tables")
	{
		tables.GET(
			"/qr/:token",
			diningTableHandler.GetDiningTableByQRToken,
		)
	}

	// Admin:
	adminTables := api.Group("/admin/tables")

	adminTables.Use(
		middleware.AuthMiddleware(jwtManager),
	)

	{
		adminTables.GET(
			"",
			diningTableHandler.GetAllDiningTables,
		)

		adminTables.GET(
			"/:id",
			diningTableHandler.GetDiningTableByID,
		)

		adminTables.POST(
			"",
			diningTableHandler.CreateDiningTable,
		)

		adminTables.PUT(
			"/:id",
			diningTableHandler.UpdateDiningTable,
		)

		adminTables.DELETE(
			"/:id",
			diningTableHandler.DeleteDiningTable,
		)
	}
}