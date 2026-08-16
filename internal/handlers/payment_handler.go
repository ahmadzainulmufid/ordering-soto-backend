package handlers

import (
	"errors"
	"net/http"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(
	paymentService services.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) HandleNotification(c *gin.Context) {
	var payload dto.MidtransNotificationPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "payload notifikasi tidak valid",
		})
		return
	}

	err := h.paymentService.HandleNotification(
		c.Request.Context(),
		payload,
	)

	switch {
	case err == nil:
		c.JSON(http.StatusOK, gin.H{
			"message": "notifikasi berhasil diproses",
		})

	case errors.Is(err, services.ErrInvalidSignature):
		c.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
		})

	case errors.Is(err, services.ErrTransactionAlreadyProcessed):
		// Tetap balas 200 agar Midtrans tidak retry terus.
		c.JSON(http.StatusOK, gin.H{
			"message": "notifikasi sudah pernah diproses",
		})

	case errors.Is(err, services.ErrOrderNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
}