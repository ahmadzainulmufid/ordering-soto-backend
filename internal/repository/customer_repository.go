package repository

import (
	"context"
	"errors"
	"fmt"

	"SotoAyam/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrCustomerNotFound = errors.New("customer tidak ditemukan")

type CustomerRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.Customer, error)
	FindByID(ctx context.Context, id int64) (*models.Customer, error)
	Create(ctx context.Context, customer *models.Customer) error
	UpdateLastLogin(ctx context.Context, id int64) error
}

type customerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) CustomerRepository {
	return &customerRepository{db: db}
}

const customerColumns = `
	id,
	full_name,
	email,
	phone,
	password_hash,
	is_active,
	last_login_at,
	created_at,
	updated_at
`

func scanCustomer(row pgx.Row, c *models.Customer) error {
	return row.Scan(
		&c.ID,
		&c.FullName,
		&c.Email,
		&c.Phone,
		&c.PasswordHash,
		&c.IsActive,
		&c.LastLoginAt,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
}

func (r *customerRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.Customer, error) {
	query := `SELECT ` + customerColumns + `
		FROM customers
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1`

	var customer models.Customer

	err := scanCustomer(r.db.QueryRow(ctx, query, email), &customer)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCustomerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mencari customer berdasarkan email: %w", err)
	}

	return &customer, nil
}

func (r *customerRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.Customer, error) {
	query := `SELECT ` + customerColumns + `
		FROM customers
		WHERE id = $1
		LIMIT 1`

	var customer models.Customer

	err := scanCustomer(r.db.QueryRow(ctx, query, id), &customer)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCustomerNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mencari customer berdasarkan id: %w", err)
	}

	return &customer, nil
}

func (r *customerRepository) Create(
	ctx context.Context,
	customer *models.Customer,
) error {
	const query = `
		INSERT INTO customers (
			full_name,
			email,
			phone,
			password_hash,
			is_active
		)
		VALUES ($1, LOWER($2), $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		customer.FullName,
		customer.Email,
		customer.Phone,
		customer.PasswordHash,
		customer.IsActive,
	).Scan(
		&customer.ID,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		// 23505 = unique_violation (email sudah dipakai)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyInUse
		}
		return fmt.Errorf("gagal membuat customer: %w", err)
	}

	return nil
}

func (r *customerRepository) UpdateLastLogin(
	ctx context.Context,
	id int64,
) error {
	const query = `
		UPDATE customers
		SET
			last_login_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal memperbarui waktu login customer: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrCustomerNotFound
	}

	return nil
}