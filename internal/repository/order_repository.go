package repository

import (
	"context"
	"errors"
	"fmt"

	"SotoAyam/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrOrderNotFound = errors.New("order tidak ditemukan")

type OrderRepository interface {
	Create(
		ctx context.Context,
		tx pgx.Tx,
		order *models.Order,
	) error

	CreateItem(
		ctx context.Context,
		tx pgx.Tx,
		item *models.OrderItem,
	) error

	FindByID(
		ctx context.Context,
		id int64,
	) (*models.Order, error)

	FindByCode(
		ctx context.Context,
		orderCode string,
	) (*models.Order, error)

	FindItemsByOrderID(
		ctx context.Context,
		orderID int64,
	) ([]models.OrderItem, error)

	FindAll(
		ctx context.Context,
	) ([]models.Order, error)

	UpdateStatus(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
		status string,
	) error

	CreateStatusHistory(
		ctx context.Context,
		tx pgx.Tx,
		history *models.OrderStatusHistory,
	) error

	BeginTx(
		ctx context.Context,
	) (pgx.Tx, error)

	UpdateSnapToken(
	ctx context.Context,
	orderID int64,
	token string,
	) error

	UpdatePaymentStatus(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
		paymentStatus string,
	) error
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(
	db *pgxpool.Pool,
) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) BeginTx(
	ctx context.Context,
) (pgx.Tx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal memulai transaction: %w",
			err,
		)
	}

	return tx, nil
}

func (r *orderRepository) Create(
	ctx context.Context,
	tx pgx.Tx,
	order *models.Order,
) error {
	const query = `
		INSERT INTO orders (
			order_code,
			customer_name,
			customer_phone,
			cashier_id,
			table_id,
			order_type,
			status,
			payment_method,
			payment_status,
			delivery_address,
			notes,
			subtotal,
			delivery_fee,
			discount,
			total
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15
		)
		RETURNING id, created_at, updated_at
	`

	err := tx.QueryRow(
		ctx,
		query,
		order.OrderCode,
		order.CustomerName,
		order.CustomerPhone,
		order.CashierID,
		order.TableID,
		order.OrderType,
		order.Status,
		order.PaymentMethod,
		order.PaymentStatus,
		order.DeliveryAddress,
		order.Notes,
		order.Subtotal,
		order.DeliveryFee,
		order.Discount,
		order.Total,
	).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal membuat order: %w",
			err,
		)
	}

	return nil
}

func (r *orderRepository) CreateItem(
	ctx context.Context,
	tx pgx.Tx,
	item *models.OrderItem,
) error {
	const query = `
		INSERT INTO order_items (
			order_id,
			product_id,
			product_name,
			product_price,
			quantity,
			subtotal,
			notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := tx.QueryRow(
		ctx,
		query,
		item.OrderID,
		item.ProductID,
		item.ProductName,
		item.ProductPrice,
		item.Quantity,
		item.Subtotal,
		item.Notes,
	).Scan(
		&item.ID,
		&item.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal membuat order item: %w",
			err,
		)
	}

	return nil
}

func (r *orderRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.Order, error) {
	const query = `
		SELECT
			id,
			order_code,
			customer_name,
			customer_phone,
			cashier_id,
			table_id,
			order_type,
			status,
			payment_method,
			payment_status,
			delivery_address,
			notes,
			subtotal,
			delivery_fee,
			discount,
			total,
			snap_token,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		LIMIT 1
	`

	var order models.Order

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&order.ID,
		&order.OrderCode,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.CashierID,
		&order.TableID,
		&order.OrderType,
		&order.Status,
		&order.PaymentMethod,
		&order.PaymentStatus,
		&order.DeliveryAddress,
		&order.Notes,
		&order.Subtotal,
		&order.DeliveryFee,
		&order.Discount,
		&order.Total,
		&order.SnapToken,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil order: %w",
			err,
		)
	}

	return &order, nil
}

func (r *orderRepository) FindByCode(
	ctx context.Context,
	orderCode string,
) (*models.Order, error) {
	const query = `
		SELECT
			id,
			order_code,
			customer_name,
			customer_phone,
			cashier_id,
			table_id,
			order_type,
			status,
			payment_method,
			payment_status,
			delivery_address,
			notes,
			subtotal,
			delivery_fee,
			discount,
			total,
			snap_token,
			created_at,
			updated_at
		FROM orders
		WHERE order_code = $1
		LIMIT 1
	`

	var order models.Order

	err := r.db.QueryRow(
		ctx,
		query,
		orderCode,
	).Scan(
		&order.ID,
		&order.OrderCode,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.CashierID,
		&order.TableID,
		&order.OrderType,
		&order.Status,
		&order.PaymentMethod,
		&order.PaymentStatus,
		&order.DeliveryAddress,
		&order.Notes,
		&order.Subtotal,
		&order.DeliveryFee,
		&order.Discount,
		&order.Total,
		&order.SnapToken,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrderNotFound
	}

	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil order berdasarkan kode: %w",
			err,
		)
	}

	return &order, nil
}

func (r *orderRepository) FindItemsByOrderID(
	ctx context.Context,
	orderID int64,
) ([]models.OrderItem, error) {
	const query = `
		SELECT
			id,
			order_id,
			product_id,
			product_name,
			product_price,
			quantity,
			subtotal,
			notes,
			created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil item order: %w",
			err,
		)
	}
	defer rows.Close()

	items := make([]models.OrderItem, 0)

	for rows.Next() {
		var item models.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.ProductPrice,
			&item.Quantity,
			&item.Subtotal,
			&item.Notes,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"gagal membaca item order: %w",
				err,
			)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca hasil item order: %w",
			err,
		)
	}

	return items, nil
}

func (r *orderRepository) FindAll(
	ctx context.Context,
) ([]models.Order, error) {
	const query = `
		SELECT
			id,
			order_code,
			customer_name,
			customer_phone,
			cashier_id,
			table_id,
			order_type,
			status,
			payment_method,
			payment_status,
			delivery_address,
			notes,
			subtotal,
			delivery_fee,
			discount,
			total,
			created_at,
			updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal mengambil daftar order: %w",
			err,
		)
	}
	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		var order models.Order

		if err := rows.Scan(
			&order.ID,
			&order.OrderCode,
			&order.CustomerName,
			&order.CustomerPhone,
			&order.CashierID,
			&order.TableID,
			&order.OrderType,
			&order.Status,
			&order.PaymentMethod,
			&order.PaymentStatus,
			&order.DeliveryAddress,
			&order.Notes,
			&order.Subtotal,
			&order.DeliveryFee,
			&order.Discount,
			&order.Total,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"gagal membaca data order: %w",
				err,
			)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"gagal membaca hasil order: %w",
			err,
		)
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	status string,
) error {
	const query = `
		UPDATE orders
		SET
			status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := tx.Exec(
		ctx,
		query,
		status,
		orderID,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal memperbarui status order: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (r *orderRepository) CreateStatusHistory(
	ctx context.Context,
	tx pgx.Tx,
	history *models.OrderStatusHistory,
) error {
	const query = `
		INSERT INTO order_status_histories (
			order_id,
			changed_by,
			status,
			description
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := tx.QueryRow(
		ctx,
		query,
		history.OrderID,
		history.ChangedBy,
		history.Status,
		history.Description,
	).Scan(
		&history.ID,
		&history.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf(
			"gagal membuat history status order: %w",
			err,
		)
	}

	return nil
}

func (r *orderRepository) UpdateSnapToken(
	ctx context.Context,
	orderID int64,
	token string,
) error {
	const query = `
		UPDATE orders
		SET
			snap_token = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		token,
		orderID,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal menyimpan snap token: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (r *orderRepository) UpdatePaymentStatus(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	paymentStatus string,
) error {
	const query = `
		UPDATE orders
		SET
			payment_status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := tx.Exec(
		ctx,
		query,
		paymentStatus,
		orderID,
	)
	if err != nil {
		return fmt.Errorf(
			"gagal memperbarui payment status: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}