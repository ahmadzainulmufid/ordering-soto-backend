package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Proses request
		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Printf(
			"[HTTP] %s %s | Status: %d | IP: %s | Latency: %s",
			method,
			path,
			statusCode,
			clientIP,
			latency,
		)

		// Catat error jika ada
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Printf(
					"[ERROR] %s %s | %v",
					method,
					path,
					err.Err,
				)
			}
		}
	}
}