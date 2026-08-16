package routes

import (
	"SotoAyam/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(
	router *gin.RouterGroup,
	paymentHandler *handlers.PaymentHandler,
) {
	router.POST(
		"/payments/notification",
		paymentHandler.HandleNotification,
	)
}