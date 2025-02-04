package constants

// OrderStatus represents the current state of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
	OrderStatusRefunded   OrderStatus = "REFUNDED"
)

// IsValid checks if the order status is valid
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusProcessing,
		OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled,
		OrderStatusRefunded:
		return true
	}
	return false
}

// String returns the string representation of the OrderStatus
func (s OrderStatus) String() string {
	return string(s)
}

// PaymentStatus represents the current state of a payment
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusFailed   PaymentStatus = "FAILED"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

// IsValid checks if the payment status is valid
func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusPending, PaymentStatusPaid,
		PaymentStatusFailed, PaymentStatusRefunded:
		return true
	}
	return false
}

// String returns the string representation of the PaymentStatus
func (s PaymentStatus) String() string {
	return string(s)
}
