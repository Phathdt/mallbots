package interfaces

import (
	"context"
	"mallbots/modules/order/application/dto"

	"github.com/phathdt/service-context/core"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID int32, req *dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetOrder(ctx context.Context, orderID int32) (*dto.OrderResponse, error)
	GetUserOrders(ctx context.Context, userID int32, paging *core.Paging) ([]*dto.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, orderID int32, req *dto.UpdateOrderStatusRequest) error
	UpdatePaymentStatus(ctx context.Context, orderID int32, req *dto.UpdatePaymentStatusRequest) error
}
