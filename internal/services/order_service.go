package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"SotoAyam/internal/dto"
	"SotoAyam/internal/models"
	"SotoAyam/internal/repository"
)

var (
	ErrOrderNotFound       = errors.New("order tidak ditemukan")
	ErrInvalidOrderType    = errors.New("tipe order tidak valid")
	ErrDiningTableRequired = errors.New("meja wajib dipilih untuk dine in")
	ErrDiningTableInactive = errors.New("meja tidak aktif")
	ErrDeliveryAddress     = errors.New("alamat pengiriman wajib diisi")
	ErrProductUnavailable  = errors.New("produk tidak tersedia")
	ErrInsufficientStock   = errors.New("stok produk tidak mencukupi")
	ErrInvalidOrderStatus  = errors.New("perubahan status order tidak valid")
)

type OrderService interface {
	GetAllOrders(
		ctx context.Context,
	) ([]dto.OrderResponse, error)

	GetOrderByID(
		ctx context.Context,
		orderID int64,
	) (*dto.OrderResponse, error)

	GetOrderByCode(
		ctx context.Context,
		orderCode string,
	) (*dto.OrderResponse, error)

	CreateOrder(
		ctx context.Context,
		request dto.CreateOrderRequest,
	) (*dto.OrderResponse, error)

	UpdateOrderStatus(
		ctx context.Context,
		orderID int64,
		status string,
		changedBy *int64,
	) (*dto.OrderResponse, error)
}

type orderService struct {
	orderRepository       repository.OrderRepository
	productRepository     repository.ProductRepository
	diningTableRepository repository.DiningTableRepository
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	productRepository repository.ProductRepository,
	diningTableRepository repository.DiningTableRepository,
) OrderService {
	return &orderService{
		orderRepository:       orderRepository,
		productRepository:     productRepository,
		diningTableRepository: diningTableRepository,
	}
}

func (s *orderService) GetAllOrders(
	ctx context.Context,
) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get orders: %w",
			err,
		)
	}

	responses := make(
		[]dto.OrderResponse,
		0,
		len(orders),
	)

	for _, order := range orders {
		items, err :=
			s.orderRepository.FindItemsByOrderID(
				ctx,
				order.ID,
			)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to get order items: %w",
				err,
			)
		}

		response := mapOrderResponse(
			&order,
			items,
		)

		responses = append(
			responses,
			response,
		)
	}

	return responses, nil
}

func (s *orderService) GetOrderByID(
	ctx context.Context,
	orderID int64,
) (*dto.OrderResponse, error) {
	order, err := s.orderRepository.FindByID(
		ctx,
		orderID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrOrderNotFound,
		) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf(
			"failed to get order: %w",
			err,
		)
	}

	items, err :=
		s.orderRepository.FindItemsByOrderID(
			ctx,
			order.ID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get order items: %w",
			err,
		)
	}

	response := mapOrderResponse(
		order,
		items,
	)

	return &response, nil
}

func (s *orderService) GetOrderByCode(
	ctx context.Context,
	orderCode string,
) (*dto.OrderResponse, error) {
	orderCode = strings.TrimSpace(orderCode)

	order, err := s.orderRepository.FindByCode(
		ctx,
		orderCode,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrOrderNotFound,
		) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf(
			"failed to get order by code: %w",
			err,
		)
	}

	items, err :=
		s.orderRepository.FindItemsByOrderID(
			ctx,
			order.ID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get order items: %w",
			err,
		)
	}

	response := mapOrderResponse(
		order,
		items,
	)

	return &response, nil
}

func (s *orderService) CreateOrder(
	ctx context.Context,
	request dto.CreateOrderRequest,
) (*dto.OrderResponse, error) {
	customerName := strings.TrimSpace(
		request.CustomerName,
	)

	orderType := strings.ToLower(
		strings.TrimSpace(request.OrderType),
	)

	if orderType != "dine_in" &&
		orderType != "takeaway" &&
		orderType != "delivery" {
		return nil, ErrInvalidOrderType
	}

	// Validasi Dine In

	var tableID sql.NullInt64

	if orderType == "dine_in" {
		if request.TableID == nil {
			return nil, ErrDiningTableRequired
		}

		table, err :=
			s.diningTableRepository.FindByID(
				ctx,
				*request.TableID,
			)

		if err != nil {
			if errors.Is(
				err,
				repository.ErrDiningTableNotFound,
			) {
				return nil, ErrDiningTableNotFound
			}

			return nil, fmt.Errorf(
				"failed to check dining table: %w",
				err,
			)
		}

		if !table.IsActive {
			return nil, ErrDiningTableInactive
		}

		tableID = sql.NullInt64{
			Int64: table.ID,
			Valid: true,
		}
	}

	// Validasi Delivery

	deliveryAddress :=
		strings.TrimSpace(
			request.DeliveryAddress,
		)

	if orderType == "delivery" &&
		deliveryAddress == "" {
		return nil, ErrDeliveryAddress
	}

	// Validasi Items

	if len(request.Items) == 0 {
		return nil, errors.New(
			"order harus memiliki minimal satu produk",
		)
	}

	orderItems := make(
		[]models.OrderItem,
		0,
		len(request.Items),
	)

	var subtotal float64

	for _, requestedItem := range request.Items {
		if requestedItem.Quantity <= 0 {
			return nil, errors.New(
				"jumlah produk harus lebih dari 0",
			)
		}

		product, err :=
			s.productRepository.FindByID(
				ctx,
				requestedItem.ProductID,
			)

		if err != nil {
			if errors.Is(
				err,
				repository.ErrProductNotFound,
			) {
				return nil, fmt.Errorf(
					"%w: ID %d",
					ErrProductNotFound,
					requestedItem.ProductID,
				)
			}

			return nil, fmt.Errorf(
				"failed to get product: %w",
				err,
			)
		}

		if !product.IsAvailable {
			return nil, fmt.Errorf(
				"%w: %s",
				ErrProductUnavailable,
				product.Name,
			)
		}

		if product.Stock <
			requestedItem.Quantity {
			return nil, fmt.Errorf(
				"%w untuk produk %s",
				ErrInsufficientStock,
				product.Name,
			)
		}

		itemSubtotal :=
			product.Price *
				float64(requestedItem.Quantity)

		item := models.OrderItem{
			ProductID: sql.NullInt64{
				Int64: product.ID,
				Valid: true,
			},
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     requestedItem.Quantity,
			Subtotal:     itemSubtotal,
		}

		if strings.TrimSpace(
			requestedItem.Notes,
		) != "" {
			item.Notes = sql.NullString{
				String: strings.TrimSpace(
					requestedItem.Notes,
				),
				Valid: true,
			}
		}

		orderItems = append(
			orderItems,
			item,
		)

		subtotal += itemSubtotal
	}

	// Biaya Order

	var deliveryFee float64

	if orderType == "delivery" {
		// Untuk sekarang fixed.
		// Nanti bisa dipindahkan ke settings/database.
		deliveryFee = 5000
	}

	discount := float64(0)

	total :=
		subtotal +
			deliveryFee -
			discount

	// Buat Order

	order := &models.Order{
		OrderCode:     generateOrderCode(),
		CustomerName:  customerName,
		TableID:       tableID,
		OrderType:     orderType,
		Status:        "pending",
		PaymentStatus: "unpaid",
		Subtotal:      subtotal,
		DeliveryFee:   deliveryFee,
		Discount:      discount,
		Total:         total,
	}

	if strings.TrimSpace(
		request.CustomerPhone,
	) != "" {
		order.CustomerPhone = sql.NullString{
			String: strings.TrimSpace(
				request.CustomerPhone,
			),
			Valid: true,
		}
	}

	if deliveryAddress != "" {
		order.DeliveryAddress = sql.NullString{
			String: deliveryAddress,
			Valid:  true,
		}
	}

	if strings.TrimSpace(
		request.Notes,
	) != "" {
		order.Notes = sql.NullString{
			String: strings.TrimSpace(
				request.Notes,
			),
			Valid: true,
		}
	}

	// Database Transaction

	tx, err :=
		s.orderRepository.BeginTx(ctx)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to start transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	if err := s.orderRepository.Create(
		ctx,
		tx,
		order,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to create order: %w",
			err,
		)
	}

	// Insert Order Items

	for i := range orderItems {
		orderItems[i].OrderID = order.ID

		if err :=
			s.orderRepository.CreateItem(
				ctx,
				tx,
				&orderItems[i],
			); err != nil {
			return nil, fmt.Errorf(
				"failed to create order item: %w",
				err,
			)
		}
	}

	// Status History

	history := &models.OrderStatusHistory{
		OrderID: order.ID,
		Status:  "pending",
		Description: sql.NullString{
			String: "Order dibuat",
			Valid:  true,
		},
	}

	if err :=
		s.orderRepository.CreateStatusHistory(
			ctx,
			tx,
			history,
		); err != nil {
		return nil, fmt.Errorf(
			"failed to create order status history: %w",
			err,
		)
	}

	// Commit

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"failed to commit transaction: %w",
			err,
		)
	}

	response := mapOrderResponse(
		order,
		orderItems,
	)

	return &response, nil
}

func (s *orderService) UpdateOrderStatus(
	ctx context.Context,
	orderID int64,
	newStatus string,
	changedBy *int64,
) (*dto.OrderResponse, error) {
	order, err :=
		s.orderRepository.FindByID(
			ctx,
			orderID,
		)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrOrderNotFound,
		) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf(
			"failed to get order: %w",
			err,
		)
	}

	newStatus = strings.ToLower(
		strings.TrimSpace(newStatus),
	)

	if !isValidStatusTransition(
		order.Status,
		newStatus,
	) {
		return nil, ErrInvalidOrderStatus
	}

	tx, err :=
		s.orderRepository.BeginTx(ctx)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to start transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	if err :=
		s.orderRepository.UpdateStatus(
			ctx,
			tx,
			orderID,
			newStatus,
		); err != nil {
		if errors.Is(
			err,
			repository.ErrOrderNotFound,
		) {
			return nil, ErrOrderNotFound
		}

		return nil, fmt.Errorf(
			"failed to update order status: %w",
			err,
		)
	}

	history := &models.OrderStatusHistory{
		OrderID: orderID,
		Status:  newStatus,
		Description: sql.NullString{
			String: fmt.Sprintf(
				"Status berubah dari %s menjadi %s",
				order.Status,
				newStatus,
			),
			Valid: true,
		},
	}

	if changedBy != nil {
		history.ChangedBy = sql.NullInt64{
			Int64: *changedBy,
			Valid: true,
		}
	}

	if err :=
		s.orderRepository.CreateStatusHistory(
			ctx,
			tx,
			history,
		); err != nil {
		return nil, fmt.Errorf(
			"failed to create status history: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"failed to commit transaction: %w",
			err,
		)
	}

	order.Status = newStatus

	items, err :=
		s.orderRepository.FindItemsByOrderID(
			ctx,
			order.ID,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get order items: %w",
			err,
		)
	}

	response := mapOrderResponse(
		order,
		items,
	)

	return &response, nil
}

func isValidStatusTransition(
	currentStatus string,
	newStatus string,
) bool {
	allowedTransitions := map[string][]string{
		"pending": {
			"confirmed",
			"cancelled",
		},
		"confirmed": {
			"cooking",
			"cancelled",
		},
		"cooking": {
			"ready",
			"cancelled",
		},
		"ready": {
			"served",
			"delivering",
			"completed",
			"cancelled",
		},
		"served": {
			"completed",
		},
		"delivering": {
			"completed",
		},
	}

	allowedStatuses, exists :=
		allowedTransitions[currentStatus]

	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}

func generateOrderCode() string {
	now := time.Now()

	return fmt.Sprintf(
		"ORD-%s",
		now.Format("20060102-150405000"),
	)
}

func mapOrderResponse(
	order *models.Order,
	items []models.OrderItem,
) dto.OrderResponse {
	itemResponses := make(
		[]dto.OrderItemResponse,
		0,
		len(items),
	)

	for _, item := range items {
		var productID *int64

		if item.ProductID.Valid {
			id := item.ProductID.Int64
			productID = &id
		}

		notes := ""

		if item.Notes.Valid {
			notes = item.Notes.String
		}

		itemResponses = append(
			itemResponses,
			dto.OrderItemResponse{
				ID:           item.ID,
				ProductID:    productID,
				ProductName:  item.ProductName,
				ProductPrice: item.ProductPrice,
				Quantity:     item.Quantity,
				Subtotal:     item.Subtotal,
				Notes:        notes,
			},
		)
	}

	response := dto.OrderResponse{
		ID:            order.ID,
		OrderCode:     order.OrderCode,
		CustomerName:  order.CustomerName,
		OrderType:     order.OrderType,
		Status:        order.Status,
		PaymentStatus: order.PaymentStatus,
		Subtotal:      order.Subtotal,
		DeliveryFee:   order.DeliveryFee,
		Discount:      order.Discount,
		Total:         order.Total,
		Items:         itemResponses,
	}

	if order.CustomerPhone.Valid {
		response.CustomerPhone =
			order.CustomerPhone.String
	}

	if order.PaymentMethod.Valid {
		response.PaymentMethod =
			order.PaymentMethod.String
	}

	if order.DeliveryAddress.Valid {
		response.DeliveryAddress =
			order.DeliveryAddress.String
	}

	if order.Notes.Valid {
		response.Notes =
			order.Notes.String
	}

	return response
}