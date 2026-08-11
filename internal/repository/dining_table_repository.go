package repository

import (
	"SotoAyam/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDiningTableNotFound = errors.New("meja makan tidak ditemukan")
)

type DiningTableRepository interface {
	FindAll(
		ctx context.Context,
	) ([]models.DiningTable, error)

	FindByID(
		ctx context.Context,
		id int64,
	) (*models.DiningTable, error)

	FindByTableNumber(
		ctx context.Context,
		tableNumber string,
	) (*models.DiningTable, error)

	FindByQRToken(
		ctx context.Context,
		qrToken string,
	) (*models.DiningTable, error)

	Create(
		ctx context.Context,
		diningTable *models.DiningTable,
	) error

	Update(
		ctx context.Context,
		diningTable *models.DiningTable,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error
}

type diningTableRepository struct {
	db *pgxpool.Pool
}

func NewDiningTableRepository(
	db *pgxpool.Pool,
) DiningTableRepository {
	return &diningTableRepository{
		db: db,
	}
}

func (r *diningTableRepository) FindAll(
	ctx context.Context,
) ([]models.DiningTable, error) {
	const query = `
		SELECT
			id,
			table_number,
			qr_token,
			is_active,
			created_at,
			updated_at
		FROM dining_tables
		ORDER BY table_number ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get dining tables: %w",
			err,
		)
	}
	defer rows.Close()

	diningTables := make([]models.DiningTable, 0)

	for rows.Next() {
		var diningTable models.DiningTable

		err := rows.Scan(
			&diningTable.ID,
			&diningTable.TableNumber,
			&diningTable.QRToken,
			&diningTable.IsActive,
			&diningTable.CreatedAt,
			&diningTable.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan dining table: %w",
				err,
			)
		}

		diningTables = append(
			diningTables,
			diningTable,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to iterate dining tables: %w",
			err,
		)
	}

	return diningTables, nil
}

func (r *diningTableRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.DiningTable, error) {
	const query = `
		SELECT
			id,
			table_number,
			qr_token,
			is_active,
			created_at,
			updated_at
		FROM dining_tables
		WHERE id = $1
		LIMIT 1
	`

	var diningTable models.DiningTable

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&diningTable.ID,
		&diningTable.TableNumber,
		&diningTable.QRToken,
		&diningTable.IsActive,
		&diningTable.CreatedAt,
		&diningTable.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrDiningTableNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get dining table: %w",
			err,
		)
	}

	return &diningTable, nil
}

func (r *diningTableRepository) FindByTableNumber(
	ctx context.Context,
	tableNumber string,
) (*models.DiningTable, error) {
	const query = `
		SELECT
			id,
			table_number,
			qr_token,
			is_active,
			created_at,
			updated_at
		FROM dining_tables
		WHERE table_number = $1
		LIMIT 1
	`

	var diningTable models.DiningTable

	err := r.db.QueryRow(
		ctx,
		query,
		tableNumber,
	).Scan(
		&diningTable.ID,
		&diningTable.TableNumber,
		&diningTable.QRToken,
		&diningTable.IsActive,
		&diningTable.CreatedAt,
		&diningTable.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrDiningTableNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find dining table by number: %w",
			err,
		)
	}

	return &diningTable, nil
}

func (r *diningTableRepository) FindByQRToken(
	ctx context.Context,
	qrToken string,
) (*models.DiningTable, error) {
	const query = `
		SELECT
			id,
			table_number,
			qr_token,
			is_active,
			created_at,
			updated_at
		FROM dining_tables
		WHERE qr_token = $1
		LIMIT 1
	`

	var diningTable models.DiningTable

	err := r.db.QueryRow(
		ctx,
		query,
		qrToken,
	).Scan(
		&diningTable.ID,
		&diningTable.TableNumber,
		&diningTable.QRToken,
		&diningTable.IsActive,
		&diningTable.CreatedAt,
		&diningTable.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrDiningTableNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find dining table by QR token: %w",
			err,
		)
	}

	return &diningTable, nil
}

func (r *diningTableRepository) Create(
	ctx context.Context,
	diningTable *models.DiningTable,
) error {
	const query = `
		INSERT INTO dining_tables (
			table_number,
			qr_token,
			is_active
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		diningTable.TableNumber,
		diningTable.QRToken,
		diningTable.IsActive,
	).Scan(
		&diningTable.ID,
		&diningTable.CreatedAt,
		&diningTable.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create dining table: %w",
			err,
		)
	}

	return nil
}

func (r *diningTableRepository) Update(
	ctx context.Context,
	diningTable *models.DiningTable,
) error {
	const query = `
		UPDATE dining_tables
		SET
			table_number = $1,
			qr_token = $2,
			is_active = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		diningTable.TableNumber,
		diningTable.QRToken,
		diningTable.IsActive,
		diningTable.ID,
	).Scan(
		&diningTable.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrDiningTableNotFound
	}

	if err != nil {
		return fmt.Errorf(
			"failed to update dining table: %w",
			err,
		)
	}

	return nil
}

func (r *diningTableRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	const query = `
		DELETE FROM dining_tables
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to delete dining table: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrDiningTableNotFound
	}

	return nil
}