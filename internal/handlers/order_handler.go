package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/services"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(
	orderService services.OrderService,
) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(
	c *gin.Context,
) {
	var request dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data order tidak valid",
			err.Error(),
		)
		return
	}

	response, err :=
		h.orderService.CreateOrder(
			c.Request.Context(),
			request,
		)

	if err != nil {
		switch {
		case errors.Is(
			err,
			services.ErrInvalidOrderType,
		):
			utils.ErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrDiningTableRequired,
		):
			utils.ErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrDiningTableNotFound,
		):
			utils.ErrorResponse(
				c,
				http.StatusNotFound,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrDiningTableInactive,
		):
			utils.ErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrDeliveryAddress,
		):
			utils.ErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrProductNotFound,
		):
			utils.ErrorResponse(
				c,
				http.StatusNotFound,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrProductUnavailable,
		):
			utils.ErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrInsufficientStock,
		):
			utils.ErrorResponse(
				c,
				http.StatusConflict,
				err.Error(),
				nil,
			)

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"order berhasil dibuat",
		response,
	)
}

func (h *OrderHandler) GetAllOrders(
	c *gin.Context,
) {
	response, err :=
		h.orderService.GetAllOrders(
			c.Request.Context(),
		)

	if err != nil {
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"data order berhasil diambil",
		response,
	)
}

func (h *OrderHandler) GetOrderByCode(
	c *gin.Context,
) {
	orderCode := strings.TrimSpace(
		c.Param("code"),
	)

	if orderCode == "" {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"kode order tidak valid",
			nil,
		)
		return
	}

	response, err :=
		h.orderService.GetOrderByCode(
			c.Request.Context(),
			orderCode,
		)

	if err != nil {
		switch {
		case errors.Is(
			err,
			services.ErrOrderNotFound,
		):
			utils.ErrorResponse(
				c,
				http.StatusNotFound,
				err.Error(),
				nil,
			)

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"data order berhasil diambil",
		response,
	)
}

func (h *OrderHandler) UpdateOrderStatus(
	c *gin.Context,
) {
	id, err := strconv.ParseInt(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id <= 0 {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"id order tidak valid",
			nil,
		)
		return
	}

	var request dto.UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"status order tidak valid",
			err.Error(),
		)
		return
	}

	var changedBy *int64

	if userIDValue, exists := c.Get("user_id"); exists {
		switch value := userIDValue.(type) {
		case int64:
			changedBy = &value

		case int:
			id := int64(value)
			changedBy = &id

		case string:
			parsedID, err := strconv.ParseInt(
				value,
				10,
				64,
			)

			if err == nil {
				changedBy = &parsedID
			}
		}
	}

	response, err :=
		h.orderService.UpdateOrderStatus(
			c.Request.Context(),
			id,
			request.Status,
			changedBy,
		)

	if err != nil {
		switch {
		case errors.Is(
			err,
			services.ErrOrderNotFound,
		):
			utils.ErrorResponse(
				c,
				http.StatusNotFound,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrInvalidOrderStatus,
		):
			utils.ErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
				nil,
			)

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"status order berhasil diperbarui",
		response,
	)
}

