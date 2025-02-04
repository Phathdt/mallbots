package dto

import "time"

type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address" validate:"required"`
	ShippingCity    string `json:"shipping_city" validate:"required"`
	ShippingCountry string `json:"shipping_country" validate:"required"`
	ShippingZip     string `json:"shipping_zip" validate:"required"`
}

type OrderItemResponse struct {
	ID        int32   `json:"id"`
	ProductID int32   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderResponse struct {
	ID              int32               `json:"id"`
	Status          string              `json:"status"`
	PaymentStatus   string              `json:"payment_status"`
	TotalAmount     float64             `json:"total_amount"`
	ShippingAddress string              `json:"shipping_address"`
	ShippingCity    string              `json:"shipping_city"`
	ShippingCountry string              `json:"shipping_country"`
	ShippingZip     string              `json:"shipping_zip"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" validate:"required"`
}
