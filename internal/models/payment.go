package models

import (
	"database/sql"
	"time"
)

type PaymentTransaction struct {
	ID                int64           `db:"id"`
	OrderID           int64           `db:"order_id"`
	TransactionID     string          `db:"transaction_id"`
	TransactionStatus string          `db:"transaction_status"`
	PaymentType       sql.NullString  `db:"payment_type"`
	GrossAmount       sql.NullFloat64 `db:"gross_amount"`
	FraudStatus       sql.NullString  `db:"fraud_status"`
	RawResponse       []byte          `db:"raw_response"`
	CreatedAt         time.Time       `db:"created_at"`
}