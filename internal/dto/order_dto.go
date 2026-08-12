package dto

type CreateOrderItemRequest struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	Notes     string `json:"notes"`
}

type CreateOrderRequest struct {
	CustomerName    string                   `json:"customer_name" binding:"required"`
	CustomerPhone   string                   `json:"customer_phone"`
	OrderType       string                   `json:"order_type" binding:"required,oneof=dine_in takeaway delivery"`
	TableID         *int64                   `json:"table_id"`
	DeliveryAddress string                   `json:"delivery_address"`
	Notes           string                   `json:"notes"`
	Items           []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=confirmed cooking ready served delivering completed cancelled"`
}

type OrderItemResponse struct {
	ID           int64   `json:"id"`
	ProductID    *int64  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
	Notes        string  `json:"notes"`
}

type OrderResponse struct {
	ID              int64               `json:"id"`
	OrderCode       string              `json:"order_code"`
	CustomerName    string              `json:"customer_name"`
	CustomerPhone   string              `json:"customer_phone"`
	OrderType       string              `json:"order_type"`
	Status          string              `json:"status"`
	PaymentMethod   string              `json:"payment_method"`
	PaymentStatus   string              `json:"payment_status"`
	DeliveryAddress string              `json:"delivery_address"`
	Notes           string              `json:"notes"`
	Subtotal        float64             `json:"subtotal"`
	DeliveryFee     float64             `json:"delivery_fee"`
	Discount        float64             `json:"discount"`
	Total           float64             `json:"total"`
	Items           []OrderItemResponse `json:"items"`
}