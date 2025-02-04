package interfaces

import (
	"context"
	"mallbots/modules/cart/domain/entities"
)

type CartRepository interface {
	Create(ctx context.Context, item *entities.CartItem) (*entities.CartItem, error)
	Update(ctx context.Context, item *entities.CartItem) error
	Delete(ctx context.Context, userID, productID int32) error
	GetByUserAndProduct(ctx context.Context, userID, productID int32) (*entities.CartItem, error)
	GetByUser(ctx context.Context, userID int32) ([]*entities.CartItem, error)
}
