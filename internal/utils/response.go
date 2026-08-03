package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(
	c *gin.Context,
	statusCode int,
	message string,
	data interface{},
) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(
	c *gin.Context,
	statusCode int,
	message string,
	errors interface{},
) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func InternalServerError(c *gin.Context) {
	ErrorResponse(
		c,
		http.StatusInternalServerError,
		"terjadi kesalahan pada server",
		nil,
	)
}