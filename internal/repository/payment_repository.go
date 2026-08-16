package repository

import (
	"context"
	"errors"
	"fmt"

	"SotoAyam/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrPaymentTransactionNotFound = errors.New(
	"payment transaction tidak ditemukan",
)

type PaymentRepository interface {
	CreateTransaction(
		ctx context.Context,
		tx pgx.Tx,
		transaction *models.PaymentTransaction,
	) error

	FindLatestByTransactionID(
		ctx context.Context,
		transactionID string,
	) (*models.PaymentTransaction, error)
}

type paymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(
	db *pgxpool.Pool,
) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreateTransaction(
	ctx context.Context,
	tx pgx.Tx,
	transaction *models.PaymentTransaction,
) error {
	const query = `
		INSERT INTO payment_transactions (
			order_id,
			transaction_id,
			transaction_status,
			payment_type,
			gross_amount,
			fraud_status,
			raw_response
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := tx.QueryRow(
		ctx,
		query,
		transaction.OrderID,
		transaction.TransactionID,
		transaction.TransactionStatus,
		transaction.PaymentType,
		transaction.GrossAmount,
		transaction.FraudStatus,
		transaction.RawResponse,
	).Scan(
		&transaction.ID,
		&transaction.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal membuat payment transaction: %w",
			err,
		)
	}

	return nil
}

func (r *paymentRepository) FindLatestByTransactionID(
	ctx context.Context,
	transactionID string,
) (*models.PaymentTransaction, error) {
	const query = `
		SELECT
			id, order_id, transaction_id, transaction_status,
			payment_type, gross_amount, fraud_status,
			raw_response, created_at
		FROM payment_transactions
		WHERE transaction_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var t models.PaymentTransaction

	err := r.db.QueryRow(
		ctx,
		query,
		transactionID,
	).Scan(
		&t.ID,
		&t.OrderID,
		&t.TransactionID,
		&t.TransactionStatus,
		&t.PaymentType,
		&t.GrossAmount,
		&t.FraudStatus,
		&t.RawResponse,
		&t.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPaymentTransactionNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil payment transaction: %w",
			err,
		)
	}

	return &t, nil
}