package services

import (
	"context"
	"mallbots/modules/cart/domain/interfaces"
	"mallbots/modules/order/application/dto"
	"mallbots/modules/order/domain/constants"
	orderEntities "mallbots/modules/order/domain/entities"
	orderInterfaces "mallbots/modules/order/domain/interfaces"
	"mallbots/shared/errorx"
	"time"

	"github.com/phathdt/service-context/core"
)

type orderService struct {
	orderRepo   orderInterfaces.OrderRepository
	cartService interfaces.CartService
}

func NewOrderService(
	orderRepo orderInterfaces.OrderRepository,
	cartService interfaces.CartService,
) orderInterfaces.OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		cartService: cartService,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, userID int32, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// Get cart items
	cartItems, err := s.cartService.GetItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errorx.ErrCartEmpty
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Create order
	order := &orderEntities.Order{
		UserID:          userID,
		Status:          constants.OrderStatusPending,
		PaymentStatus:   constants.PaymentStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: req.ShippingAddress,
		ShippingCity:    req.ShippingCity,
		ShippingCountry: req.ShippingCountry,
		ShippingZip:     req.ShippingZip,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	newOrder, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	// Create order items
	var orderItems []*orderEntities.OrderItem
	for _, item := range cartItems {
		orderItems = append(orderItems, &orderEntities.OrderItem{
			OrderID:   newOrder.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	if err := s.orderRepo.CreateOrderItems(ctx, newOrder.ID, orderItems); err != nil {
		return nil, err
	}

	// Clear cart after successful order creation
	for _, item := range cartItems {
		if err := s.cartService.RemoveItem(ctx, userID, item.ProductID); err != nil {
			// Log error but don't fail the order creation
			// TODO: Implement retry mechanism for cart cleanup
			continue
		}
	}

	newOrder.Items = orderItems
	return s.convertToResponse(newOrder), nil
}

func (s *orderService) GetOrder(ctx context.Context, orderID int32) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(order), nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID int32, paging *core.Paging) ([]*dto.OrderResponse, error) {
	orders, err := s.orderRepo.GetByUserID(ctx, userID, paging)
	if err != nil {
		return nil, err
	}

	var responses []*dto.OrderResponse
	for _, order := range orders {
		responses = append(responses, s.convertToResponse(order))
	}

	return responses, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID int32, req *dto.UpdateOrderStatusRequest) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	newStatus := constants.OrderStatus(req.Status)
	if !newStatus.IsValid() {
		return errorx.ErrInvalidOrderStatus
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, newStatus) {
		return errorx.ErrInvalidStatusTransition
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, newStatus)
}

func (s *orderService) UpdatePaymentStatus(ctx context.Context, orderID int32, req *dto.UpdatePaymentStatusRequest) error {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	newStatus := constants.PaymentStatus(req.PaymentStatus)
	if !newStatus.IsValid() {
		return errorx.ErrInvalidPaymentStatus
	}

	// Validate payment status transition
	if !s.isValidPaymentStatusTransition(order.PaymentStatus, newStatus) {
		return errorx.ErrInvalidPaymentStatusTransition
	}

	return s.orderRepo.UpdatePaymentStatus(ctx, orderID, newStatus)
}

func (s *orderService) convertToResponse(order *orderEntities.Order) *dto.OrderResponse {
	var itemResponses []dto.OrderItemResponse
	for _, item := range order.Items {
		itemResponses = append(itemResponses, dto.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &dto.OrderResponse{
		ID:              order.ID,
		Status:          order.Status.String(),
		PaymentStatus:   order.PaymentStatus.String(),
		TotalAmount:     order.TotalAmount,
		ShippingAddress: order.ShippingAddress,
		ShippingCity:    order.ShippingCity,
		ShippingCountry: order.ShippingCountry,
		ShippingZip:     order.ShippingZip,
		Items:           itemResponses,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func (s *orderService) isValidStatusTransition(current, new constants.OrderStatus) bool {
	// Define valid status transitions
	validTransitions := map[constants.OrderStatus][]constants.OrderStatus{
		constants.OrderStatusPending: {
			constants.OrderStatusConfirmed,
			constants.OrderStatusCancelled,
		},
		constants.OrderStatusConfirmed: {
			constants.OrderStatusProcessing,
			constants.OrderStatusCancelled,
		},
		constants.OrderStatusProcessing: {
			constants.OrderStatusShipped,
		},
		constants.OrderStatusShipped: {
			constants.OrderStatusDelivered,
		},
		constants.OrderStatusDelivered: {
			constants.OrderStatusRefunded,
		},
	}

	// Check if the new status is allowed
	allowedStatuses, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == new {
			return true
		}
	}

	return false
}

func (s *orderService) isValidPaymentStatusTransition(current, new constants.PaymentStatus) bool {
	// Define valid payment status transitions
	validTransitions := map[constants.PaymentStatus][]constants.PaymentStatus{
		constants.PaymentStatusPending: {
			constants.PaymentStatusPaid,
			constants.PaymentStatusFailed,
		},
		constants.PaymentStatusPaid: {
			constants.PaymentStatusRefunded,
		},
		constants.PaymentStatusFailed: {
			constants.PaymentStatusPending,
		},
	}

	// Check if the new status is allowed
	allowedStatuses, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == new {
			return true
		}
	}

	return false
}
