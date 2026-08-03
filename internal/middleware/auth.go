package middleware

import (
	"net/http"
	"strings"

	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	jwtManager *utils.JWTManager,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		if authorizationHeader == "" {
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"access token diperlukan",
				nil,
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(
			authorizationHeader,
			" ",
			2,
		)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"format authorization harus Bearer token",
				nil,
			)
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(parts[1])
		if err != nil {
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"access token tidak valid atau kedaluwarsa",
				nil,
			)
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func RequireRoles(
	allowedRoles ...string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"role user tidak ditemukan",
				nil,
			)
			c.Abort()
			return
		}

		currentRole, ok := roleValue.(string)
		if !ok {
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"role user tidak valid",
				nil,
			)
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if currentRole == allowedRole {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(
			c,
			http.StatusForbidden,
			"anda tidak memiliki izin untuk mengakses resource ini",
			nil,
		)
		c.Abort()
	}
}