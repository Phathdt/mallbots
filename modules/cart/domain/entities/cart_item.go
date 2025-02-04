package entities

import "time"

type CartItem struct {
	ID        int32
	UserID    int32
	ProductID int32
	Quantity  int32
	Price     float64 // Price at the time of adding to cart
	CreatedAt time.Time
	UpdatedAt time.Time
}
