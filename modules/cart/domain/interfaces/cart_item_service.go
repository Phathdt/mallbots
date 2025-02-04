package interfaces

import (
	"context"
	"mallbots/modules/cart/application/dto"
)

type CartItemService interface {
	AddItem(ctx context.Context, userID int32, req *dto.CartItemRequest) (*dto.CartItemResponse, error)
	UpdateQuantity(ctx context.Context, userID int32, req *dto.CartItemRequest) (*dto.CartItemResponse, error)
	RemoveItem(ctx context.Context, userID, productID int32) error
	GetItems(ctx context.Context, userID int32) ([]*dto.CartItemResponse, error)
}
