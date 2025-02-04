// modules/order/domain/entities/order.go

package entities

import (
	"mallbots/modules/order/domain/constants"
	"time"
)

type Order struct {
	ID              int32
	UserID          int32
	Status          constants.OrderStatus
	PaymentStatus   constants.PaymentStatus
	TotalAmount     float64
	ShippingAddress string
	ShippingCity    string
	ShippingCountry string
	ShippingZip     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Items           []*OrderItem
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == constants.OrderStatusPending ||
		o.Status == constants.OrderStatusConfirmed
}

func (o *Order) CanBeRefunded() bool {
	return o.Status == constants.OrderStatusDelivered &&
		o.PaymentStatus == constants.PaymentStatusPaid
}
