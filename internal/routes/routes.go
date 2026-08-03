package routes

import (
	"net/http"

	"SotoAyam/internal/handlers"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	jwtManager *utils.JWTManager,
) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(
			c,
			http.StatusOK,
			"API berjalan dengan baik",
			gin.H{
				"status": "healthy",
			},
		)
	})

	apiV1 := router.Group("/api/v1")

	RegisterAuthRoutes(
		apiV1,
		authHandler,
		jwtManager,
	)

	return router
}