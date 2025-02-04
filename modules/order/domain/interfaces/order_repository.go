package interfaces

import (
	"context"
	"mallbots/modules/order/domain/constants"
	"mallbots/modules/order/domain/entities"

	"github.com/phathdt/service-context/core"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) (*entities.Order, error)
	CreateOrderItems(ctx context.Context, orderID int32, items []*entities.OrderItem) error
	GetByID(ctx context.Context, id int32) (*entities.Order, error)
	GetByUserID(ctx context.Context, userID int32, paging *core.Paging) ([]*entities.Order, error)
	UpdateStatus(ctx context.Context, id int32, status constants.OrderStatus) error
	UpdatePaymentStatus(ctx context.Context, id int32, status constants.PaymentStatus) error
}
