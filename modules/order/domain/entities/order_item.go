package entities

import "time"

type OrderItem struct {
	ID        int32
	OrderID   int32
	ProductID int32
	Quantity  int32
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
