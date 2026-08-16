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
	ErrProductNotFound = errors.New("produk tidak ditemukan")
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]models.Product, error)

	FindByID(
		ctx context.Context,
		id int64,
	) (*models.Product, error)

	FindByName(
		ctx context.Context,
		name string,
	) ([]models.Product, error)

	FindByCategoryID(
		ctx context.Context,
		categoryID int64,
	) ([]models.Product, error)

	FindExactByName(
	ctx context.Context,
	name string,
) (*models.Product, error)

	Create(
		ctx context.Context,
		product *models.Product,
	) error

	Update(
		ctx context.Context,
		product *models.Product,
	) error

	Delete(
		ctx context.Context,
		id int64,
	) error

	ReduceStock(ctx context.Context, tx pgx.Tx, productID int64, quantity int) error
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

func (r *productRepository) FindAll(
	ctx context.Context,
) ([]models.Product, error) {
	const query = `
		SELECT
			id,
			category_id,
			name,
			slug,
			description,
			price,
			image_url,
			stock,
			is_available,
			created_at,
			updated_at
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil semua produk: %w",
			err,
		)
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.ImageURL,
			&product.Stock,
			&product.IsAvailable,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data produk: %w",
				err,
			)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal melakukan iterasi data produk: %w",
			err,
		)
	}

	return products, nil
}

func (r *productRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.Product, error) {
	const query = `
		SELECT
			id,
			category_id,
			name,
			slug,
			description,
			price,
			image_url,
			stock,
			is_available,
			created_at,
			updated_at
		FROM products
		WHERE id = $1
		LIMIT 1
	`

	var product models.Product

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
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mencari produk berdasarkan ID: %w",
			err,
		)
	}

	return &product, nil
}

func (r *productRepository) FindByName(
	ctx context.Context,
	name string,
) ([]models.Product, error) {
	const query = `
		SELECT
			id,
			category_id,
			name,
			slug,
			description,
			price,
			image_url,
			stock,
			is_available,
			created_at,
			updated_at
		FROM products
		WHERE name ILIKE $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		"%"+name+"%",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mencari produk berdasarkan nama: %w",
			err,
		)
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.ImageURL,
			&product.Stock,
			&product.IsAvailable,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data produk: %w",
				err,
			)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal melakukan iterasi hasil pencarian produk: %w",
			err,
		)
	}

	return products, nil
}

func (r *productRepository) FindByCategoryID(
	ctx context.Context,
	categoryID int64,
) ([]models.Product, error) {
	const query = `
		SELECT
			id,
			category_id,
			name,
			slug,
			description,
			price,
			image_url,
			stock,
			is_available,
			created_at,
			updated_at
		FROM products
		WHERE category_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		categoryID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil produk berdasarkan kategori: %w",
			err,
		)
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.ImageURL,
			&product.Stock,
			&product.IsAvailable,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data produk: %w",
				err,
			)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal melakukan iterasi produk berdasarkan kategori: %w",
			err,
		)
	}

	return products, nil
}

func (r *productRepository) FindExactByName(
	ctx context.Context,
	name string,
) (*models.Product, error) {
	const query = `
		SELECT
			id,
			category_id,
			name,
			slug,
			description,
			price,
			image_url,
			stock,
			is_available,
			created_at,
			updated_at
		FROM products
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	var product models.Product

	err := r.db.QueryRow(
		ctx,
		query,
		name,
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
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find product by exact name: %w",
			err,
		)
	}

	return &product, nil
}

func (r *productRepository) Create(
	ctx context.Context,
	product *models.Product,
) error {
	const query = `
		INSERT INTO products (
			category_id,
			name,
			slug,
			description,
			price,
			image_url,
			stock,
			is_available
		)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8
		)
		RETURNING
			id,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		product.CategoryID,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.ImageURL,
		product.Stock,
		product.IsAvailable,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal membuat produk: %w",
			err,
		)
	}

	return nil
}

func (r *productRepository) Update(
	ctx context.Context,
	product *models.Product,
) error {
	const query = `
		UPDATE products
		SET
			category_id = $1,
			name = $2,
			slug = $3,
			description = $4,
			price = $5,
			image_url = $6,
			stock = $7,
			is_available = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		product.CategoryID,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.ImageURL,
		product.Stock,
		product.IsAvailable,
		product.ID,
	).Scan(
		&product.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrProductNotFound
	}

	if err != nil {
		return fmt.Errorf(
			"gagal memperbarui produk: %w",
			err,
		)
	}

	return nil
}

func (r *productRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	const query = `
		DELETE FROM products
		WHERE id = $1
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal menghapus produk: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (r *productRepository) ReduceStock(
	ctx context.Context,
	tx pgx.Tx,
	productID int64,
	quantity int,
) error {
	const query = `
		UPDATE products
		SET stock = stock - $1,
			is_available = CASE WHEN (stock - $1) <= 0 THEN false ELSE is_available END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND stock >= $1
	`
	result, err := tx.Exec(ctx, query, quantity, productID)
	if err != nil {
		return fmt.Errorf("gagal mengurangi stok produk: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("stok produk tidak mencukupi atau produk tidak ditemukan")
	}
	return nil
}