package dto

import "time"

type ProductResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Price       float64   `json:"price"`
	CategoryID  int32     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductListRequest struct {
	Search   string   `query:"search"`
	MinPrice *float64 `query:"min_price"`
	MaxPrice *float64 `query:"max_price"`
	Category *int32   `query:"category"`
	SortBy   string   `query:"sort_by"`
}
