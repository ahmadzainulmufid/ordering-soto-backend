package models

import (
	"database/sql"
	"time"
)

type OrderItem struct {
	ID           int64          `db:"id"`
	OrderID      int64          `db:"order_id"`
	ProductID    sql.NullInt64  `db:"product_id"`
	ProductName  string         `db:"product_name"`
	ProductPrice float64        `db:"product_price"`
	Quantity     int            `db:"quantity"`
	Subtotal     float64        `db:"subtotal"`
	Notes        sql.NullString `db:"notes"`
	CreatedAt    time.Time      `db:"created_at"`
}