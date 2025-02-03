package entities

import "time"

type Product struct {
	ID          int32
	Name        string
	Description *string
	Price       float64
	CategoryID  int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Category struct {
	ID        int32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
