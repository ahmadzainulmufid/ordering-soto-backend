package models

import (
	"time"
)

type DiningTable struct {
	ID		   	int64          `db:"id"`
	TableNumber string         `db:"table_number"`
	QRToken	 	string         `db:"qr_token"`
	IsActive     bool           `db:"is_active"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}