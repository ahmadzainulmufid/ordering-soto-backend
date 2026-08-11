package models

import (
	"database/sql"
	"time"
)

type OrderStatusHistory struct {
	ID          int64          `db:"id"`
	OrderID     int64          `db:"order_id"`
	ChangedBy   sql.NullInt64  `db:"changed_by"`
	Status      string         `db:"status"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
}