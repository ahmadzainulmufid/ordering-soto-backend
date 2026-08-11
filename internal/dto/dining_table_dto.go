package dto

type CreateDiningTableRequest struct {
	TableNumber string `json:"table_number" binding:"required"`
}

type UpdateDiningTableRequest struct {
	TableNumber string `json:"table_number" binding:"required"`
	IsActive    bool   `json:"is_active"`
}

type DiningTableResponse struct {
	ID          int64  `json:"id"`
	TableNumber string `json:"table_number"`
	QRToken     string `json:"qr_token"`
	IsActive    bool   `json:"is_active"`
}