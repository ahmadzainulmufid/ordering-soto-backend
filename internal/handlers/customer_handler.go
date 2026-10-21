package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/repository"
	"SotoAyam/internal/services"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

type CustomerAuthHandler struct {
	customerService services.CustomerAuthService
}

func NewCustomerAuthHandler(
	customerService services.CustomerAuthService,
) *CustomerAuthHandler {
	return &CustomerAuthHandler{customerService: customerService}
}

func (h *CustomerAuthHandler) Register(c *gin.Context) {
	var request dto.CustomerRegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data registrasi tidak valid",
			err.Error(),
		)
		return
	}

	response, err := h.customerService.Register(c.Request.Context(), request)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExist) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"registrasi berhasil",
		response,
	)
}

func (h *CustomerAuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data login tidak valid",
			err.Error(),
		)
		return
	}

	response, err := h.customerService.Login(c.Request.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		case errors.Is(err, services.ErrInactiveUser):
			utils.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
		default:
			utils.InternalServerError(c)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "login berhasil", response)
}

func (h *CustomerAuthHandler) Profile(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"identitas user tidak ditemukan",
			nil,
		)
		return
	}

	customerID, err := strconv.ParseInt(userIDValue.(string), 10, 64)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"identitas user tidak valid",
			nil,
		)
		return
	}

	response, err := h.customerService.GetProfile(
		c.Request.Context(),
		customerID,
	)
	if err != nil {
		if errors.Is(err, repository.ErrCustomerNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"profil customer berhasil diambil",
		response,
	)
}