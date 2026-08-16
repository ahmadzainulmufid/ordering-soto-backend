package models

import (
	"database/sql"
	"time"
)

type Order struct {
	ID              int64          `db:"id"`
	OrderCode       string         `db:"order_code"`
	CustomerName    string         `db:"customer_name"`
	CustomerPhone   sql.NullString `db:"customer_phone"`
	CashierID       sql.NullInt64  `db:"cashier_id"`
	TableID         sql.NullInt64  `db:"table_id"`
	OrderType       string         `db:"order_type"`
	Status          string         `db:"status"`
	PaymentMethod   sql.NullString `db:"payment_method"`
	PaymentStatus   string         `db:"payment_status"`
	DeliveryAddress sql.NullString `db:"delivery_address"`
	Notes           sql.NullString `db:"notes"`
	Subtotal        float64        `db:"subtotal"`
	DeliveryFee     float64        `db:"delivery_fee"`
	Discount        float64        `db:"discount"`
	Total           float64        `db:"total"`
	SnapToken		sql.NullString	`db:"snap_token"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}