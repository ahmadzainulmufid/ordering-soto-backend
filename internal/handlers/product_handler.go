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

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(
	productService services.ProductService,
) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateProduct(
	c *gin.Context,
) {
	var request dto.CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data produk tidak valid",
			err.Error(),
		)
		return
	}

	response, err :=
		h.productService.CreateProduct(
			c.Request.Context(),
			request,
		)

	if err != nil {
		switch {
		case errors.Is(
			err,
			services.ErrProductExists,
		):
			utils.ErrorResponse(
				c,
				http.StatusConflict,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrCategoryInvalid,
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
		http.StatusCreated,
		"produk berhasil dibuat",
		response,
	)
}

func (h *ProductHandler) GetAllProducts(
	c *gin.Context,
) {
	response, err :=
		h.productService.GetAllProducts(
			c.Request.Context(),
		)

	if err != nil {
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"data produk berhasil diambil",
		response,
	)
}

func (h *ProductHandler) GetProductByID(
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
			"id produk tidak valid",
			nil,
		)
		return
	}

	response, err :=
		h.productService.GetProductByID(
			c.Request.Context(),
			id,
		)

	if err != nil {
		switch {
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

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"data produk berhasil diambil",
		response,
	)
}

func (h *ProductHandler) SearchProducts(
	c *gin.Context,
) {
	name := strings.TrimSpace(
		c.Query("q"),
	)

	if name == "" {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"kata kunci pencarian wajib diisi",
			nil,
		)
		return
	}

	response, err :=
		h.productService.SearchProducts(
			c.Request.Context(),
			name,
		)

	if err != nil {
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"hasil pencarian produk berhasil diambil",
		response,
	)
}

func (h *ProductHandler) UpdateProduct(
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
			"id produk tidak valid",
			nil,
		)
		return
	}

	var request dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"data produk tidak valid",
			err.Error(),
		)
		return
	}

	response, err :=
		h.productService.UpdateProduct(
			c.Request.Context(),
			id,
			request,
		)

	if err != nil {
		switch {
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
			services.ErrProductExists,
		):
			utils.ErrorResponse(
				c,
				http.StatusConflict,
				err.Error(),
				nil,
			)

		case errors.Is(
			err,
			services.ErrCategoryInvalid,
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
		"produk berhasil diperbarui",
		response,
	)
}

func (h *ProductHandler) DeleteProduct(
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
			"id produk tidak valid",
			nil,
		)
		return
	}

	err = h.productService.DeleteProduct(
		c.Request.Context(),
		id,
	)

	if err != nil {
		switch {
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

		default:
			utils.InternalServerError(c)
		}

		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"produk berhasil dihapus",
		nil,
	)
}