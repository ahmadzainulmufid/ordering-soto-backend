package handlers

import (
	"log"
	"net/http"

	"SotoAyam/internal/services"
	"SotoAyam/internal/utils"

	"github.com/gin-gonic/gin"
)

const maxImageSize = 2 << 20 // 2 MB

type UploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(
	uploadService services.UploadService,
) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	header, err := c.FormFile("image")
	if err != nil {
		utils.ErrorResponse(
			c, http.StatusBadRequest,
			"file gambar wajib diisi", nil,
		)
		return
	}

	if header.Size > maxImageSize {
		utils.ErrorResponse(
			c, http.StatusBadRequest,
			"ukuran gambar maksimal 2 MB", nil,
		)
		return
	}

	file, err := header.Open()
	if err != nil {
		utils.InternalServerError(c)
		return
	}
	defer file.Close()

	// cek tipe file asli dari isinya, bukan dari nama file
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])

	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
	default:
		utils.ErrorResponse(
			c, http.StatusBadRequest,
			"format gambar harus JPG, PNG, atau WEBP", nil,
		)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		utils.InternalServerError(c)
		return
	}

	url, err := h.uploadService.UploadImage(c.Request.Context(), file)
	if err != nil {
		log.Printf("upload gambar gagal: %v", err)
		utils.InternalServerError(c)
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"gambar berhasil diunggah",
		gin.H{"image_url": url},
	)
}