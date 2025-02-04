package dto

type CartItemRequest struct {
	ProductID int32 `json:"product_id" validate:"required"`
	Quantity  int32 `json:"quantity" validate:"required,min=1"`
}

type CartItemResponse struct {
	ID        int32   `json:"id"`
	ProductID int32   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}
