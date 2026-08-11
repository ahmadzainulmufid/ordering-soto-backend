package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/services"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

type DiningTableHandler struct {
	diningTableService services.DiningTableService
}

func NewDiningTableHandler(
	diningTableService services.DiningTableService,
) *DiningTableHandler {
	return &DiningTableHandler{
		diningTableService: diningTableService,
	}
}

func (h *DiningTableHandler) CreateDiningTable(
	c *gin.Context,
) {
	var request dto.CreateDiningTableRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data meja tidak valid",
			err.Error(),
		)
		return
	}

	response, err :=
		h.diningTableService.CreateDiningTable(
			c.Request.Context(),
			request,
		)

	if err != nil {
		switch {
		case errors.Is(
			err,
			services.ErrDiningTableExists,
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
		"meja berhasil dibuat",
		response,
	)
}

func (h *DiningTableHandler) GetAllDiningTables(
	c *gin.Context,
) {
	response, err :=
		h.diningTableService.GetAllDiningTables(
			c.Request.Context(),
		)

	if err != nil {
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"data meja berhasil diambil",
		response,
	)
}

func (h *DiningTableHandler) GetDiningTableByID(
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
			"id meja tidak valid",
			nil,
		)
		return
	}

	response, err :=
		h.diningTableService.GetDiningTableByID(
			c.Request.Context(),
			id,
		)

	if err != nil {
		switch {
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

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"data meja berhasil diambil",
		response,
	)
}

func (h *DiningTableHandler) UpdateDiningTable(
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
			"id meja tidak valid",
			nil,
		)
		return
	}

	var request dto.UpdateDiningTableRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data meja tidak valid",
			err.Error(),
		)
		return
	}

	response, err :=
		h.diningTableService.UpdateDiningTable(
			c.Request.Context(),
			id,
			request,
		)

	if err != nil {
		switch {
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
			services.ErrDiningTableExists,
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
		http.StatusOK,
		"meja berhasil diperbarui",
		response,
	)
}

func (h *DiningTableHandler) DeleteDiningTable(
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
			"id meja tidak valid",
			nil,
		)
		return
	}

	err = h.diningTableService.DeleteDiningTable(
		c.Request.Context(),
		id,
	)

	if err != nil {
		switch {
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

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"meja berhasil dihapus",
		nil,
	)
}

func (h *DiningTableHandler) GetDiningTableByQRToken(
	c *gin.Context,
) {
	token := c.Param("token")

	if token == "" {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"QR token tidak valid",
			nil,
		)
		return
	}

	response, err :=
		h.diningTableService.GetDiningTableByQRToken(
			c.Request.Context(),
			token,
		)

	if err != nil {
		switch {
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

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"meja berhasil ditemukan",
		response,
	)
}