package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(
	api *gin.RouterGroup,
	categoryHandler *handlers.CategoryHandler,
	jwtManager *utils.JWTManager,
) {
	categories := api.Group("/categories")
	{
		categories.GET(
			"",
			categoryHandler.GetAllCategories,
		)

		categories.GET(
			"/:id",
			categoryHandler.GetCategoryByID,
		)
	}

	adminCategories := api.Group("/admin/categories")

	adminCategories.Use(
		middleware.AuthMiddleware(jwtManager),
	)

	{
		adminCategories.POST(
			"",
			categoryHandler.CreateCategory,
		)

		adminCategories.PUT(
			"/:id",
			categoryHandler.UpdateCategory,
		)

		adminCategories.DELETE(
			"/:id",
			categoryHandler.DeleteCategory,
		)
	}
}