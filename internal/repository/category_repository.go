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
	ErrCategoryNotFound = errors.New("kategori tidak ditemukan")
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]models.Category, error)
	FindByID(ctx context.Context, id int64) (*models.Category, error)
	FindByName(ctx context.Context, name string) (*models.Category, error)

	Create(ctx context.Context, category *models.Category) error
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(
	db *pgxpool.Pool,
) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) FindAll(
	ctx context.Context,
) ([]models.Category, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			is_active,
			created_at,
			updated_at
		FROM categories
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get categories: %w",
			err,
		)
	}
	defer rows.Close()

	categories := make([]models.Category, 0)

	for rows.Next() {
		var category models.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan category: %w",
				err,
			)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to iterate categories: %w",
			err,
		)
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.Category, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			is_active,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
		LIMIT 1
	`

	var category models.Category

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCategoryNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find category by ID: %w",
			err,
		)
	}

	return &category, nil
}

func (r *categoryRepository) FindByName(
	ctx context.Context,
	name string,
) (*models.Category, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			is_active,
			created_at,
			updated_at
		FROM categories
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	var category models.Category

	err := r.db.QueryRow(
		ctx,
		query,
		name,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCategoryNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find category by name: %w",
			err,
		)
	}

	return &category, nil
}

func (r *categoryRepository) Create(
	ctx context.Context,
	category *models.Category,
) error {
	const query = `
		INSERT INTO categories (
			name,
			slug,
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
		category.Name,
		category.Slug,
		category.IsActive,
	).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create category: %w",
			err,
		)
	}

	return nil
}

func (r *categoryRepository) Update(
	ctx context.Context,
	category *models.Category,
) error {
	const query = `
		UPDATE categories
		SET
			name = $1,
			slug = $2,
			is_active = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.IsActive,
		category.ID,
	).Scan(
		&category.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrCategoryNotFound
	}

	if err != nil {
		return fmt.Errorf(
			"failed to update category: %w",
			err,
		)
	}

	return nil
}

func (r *categoryRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	const query = `
		DELETE FROM categories
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to delete category: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return nil
}