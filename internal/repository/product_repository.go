package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProductNotFound = errors.New("produk tidak ditemukan")

type Product struct {
	ID           int64
	CategoryID   int64
	Name         string
	Slug         string
	Description  string
	Price        float64
	ImageURL     string
	Stock        int
	IsAvailable  bool
}

type ProductRepository interface {
	FindByID(
		ctx context.Context,
		id int64,
	) (*Product, error)
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(
	db *pgxpool.Pool,
) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) FindByID(
	ctx context.Context,
	id int64,
) (*Product, error) {
	const query = `
		SELECT
			id,
			category_id,
			name,
			slug,
			COALESCE(description, ''),
			price,
			COALESCE(image_url, ''),
			stock,
			is_available
		FROM products
		WHERE id = $1
	`

	var product Product

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.ImageURL,
		&product.Stock,
		&product.IsAvailable,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil data produk: %w",
			err,
		)
	}

	return &product, nil
}