package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(
	router *gin.RouterGroup,
	authHandler *handlers.AuthHandler,
	jwtManager *utils.JWTManager,
) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)

		auth.GET(
			"/profile",
			middleware.AuthMiddleware(jwtManager),
			authHandler.Profile,
		)
	}

	adminUsers := router.Group("/admin/users")
	adminUsers.Use(
		middleware.AuthMiddleware(jwtManager),
		middleware.RequireRoles("owner"),
	)
	{
		adminUsers.POST("", authHandler.CreateUser)
		adminUsers.GET("", authHandler.GetAllUsers)
		adminUsers.PUT("/:id", authHandler.UpdateUser)
		adminUsers.DELETE("/:id", authHandler.DeleteUser)
	}
}