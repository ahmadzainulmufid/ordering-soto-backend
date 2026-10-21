package routes

import (
	"SotoAyam/internal/handlers"
	"SotoAyam/internal/middleware"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(
	router *gin.RouterGroup,
	uploadHandler *handlers.UploadHandler,
	jwtManager *utils.JWTManager,
) {
	adminUpload := router.Group("/admin/upload")

	adminUpload.Use(
		middleware.AuthMiddleware(jwtManager),
		middleware.RequireRoles("owner", "admin"),
	)

	{
		adminUpload.POST(
			"/image",
			uploadHandler.UploadImage,
		)
	}
}