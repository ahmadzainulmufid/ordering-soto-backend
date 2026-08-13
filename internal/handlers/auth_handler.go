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

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(
	authService services.AuthService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
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

	response, err := h.authService.Login(
		c.Request.Context(),
		request,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			utils.ErrorResponse(
				c,
				http.StatusUnauthorized,
				err.Error(),
				nil,
			)

		case errors.Is(err, services.ErrInactiveUser):
			utils.ErrorResponse(
				c,
				http.StatusForbidden,
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
		"login berhasil",
		response,
	)
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var request dto.CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data user tidak valid",
			err.Error(),
		)
		return
	}

	creatorRole, exists := c.Get("user_role")
	if !exists {
		utils.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"identitas user tidak ditemukan",
			nil,
		)
		return
	}

	response, err := h.authService.CreateUser(
		c.Request.Context(),
		request,
		creatorRole.(string),
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbiddenRole):
			utils.ErrorResponse(
				c,
				http.StatusForbidden,
				err.Error(),
				nil,
			)

		case errors.Is(err, services.ErrEmailAlreadyExist):
			utils.ErrorResponse(
				c,
				http.StatusConflict,
				err.Error(),
				nil,
			)

		default:
			utils.ErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
				nil,
			)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"akun user berhasil dibuat",
		response,
	)
}

func (h *AuthHandler) Register(c *gin.Context) {
    var request dto.CreateUserRequest 

    if err := c.ShouldBindJSON(&request); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "data registrasi tidak valid", err.Error())
        return
    }

    response, err := h.authService.CreateUser(c.Request.Context(), request, "user")
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, "registrasi berhasil", response)
}

func (h *AuthHandler) Profile(c *gin.Context) {
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

	userID, err := strconv.ParseInt(
		userIDValue.(string),
		10,
		64,
	)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"identitas user tidak valid",
			nil,
		)
		return
	}

	response, err := h.authService.GetProfile(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.ErrorResponse(
				c,
				http.StatusNotFound,
				err.Error(),
				nil,
			)
			return
		}

		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"profil user berhasil diambil",
		response,
	)
}

func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers(c.Request.Context())
	if err != nil {
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"data user berhasil diambil",
		users,
	)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	var request struct {
		FullName string `json:"full_name" binding:"required"`
		Phone    string `json:"phone"`
		Role     string `json:"role" binding:"required,oneof=admin cashier kitchen"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "data update tidak valid", err.Error())
		return
	}

	err = h.authService.UpdateUser(c.Request.Context(), userID, request.FullName, request.Phone, request.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user berhasil diperbarui", nil)
}

// Delete User
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	err = h.authService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "user berhasil dihapus", nil)
}