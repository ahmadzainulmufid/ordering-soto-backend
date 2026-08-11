package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AuthHandler        *handlers.AuthHandler
	CategoryHandler    *handlers.CategoryHandler
	ProductHandler     *handlers.ProductHandler
	DiningTableHandler *handlers.DiningTableHandler
	OrderHandler       *handlers.OrderHandler
}

func NewRouter(
	handlers Handlers,
	jwtManager *utils.JWTManager,
) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		middleware.LoggerMiddleware(),
		middleware.CORSMiddleware(),
	)

	api := router.Group("/api/v1")

	RegisterAuthRoutes(
		api,
		handlers.AuthHandler,
		jwtManager,
	)

	RegisterCategoryRoutes(
		api,
		handlers.CategoryHandler,
		jwtManager,
	)

	RegisterProductRoutes(
		api,
		handlers.ProductHandler,
		jwtManager,
	)

	RegisterDiningTableRoutes(
		api,
		handlers.DiningTableHandler,
		jwtManager,
	)

	RegisterOrderRoutes(
		api,
		handlers.OrderHandler,
		jwtManager,
	)

	return router
}