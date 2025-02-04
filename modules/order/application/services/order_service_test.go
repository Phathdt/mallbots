package services

import (
	"context"
	cartDto "mallbots/modules/cart/application/dto"
	"mallbots/modules/order/application/dto"
	"mallbots/modules/order/domain/constants"
	"mallbots/modules/order/domain/entities"
	"mallbots/modules/order/domain/interfaces"
	"mallbots/shared/errorx"
	"testing"
	"time"

	"github.com/phathdt/service-context/core"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Order), args.Error(1)
}

func (m *MockOrderRepository) CreateOrderItems(ctx context.Context, orderID int32, items []*entities.OrderItem) error {
	args := m.Called(ctx, orderID, items)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id int32) (*entities.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUserID(ctx context.Context, userID int32, paging *core.Paging) ([]*entities.Order, error) {
	args := m.Called(ctx, userID, paging)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id int32, status constants.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdatePaymentStatus(ctx context.Context, id int32, status constants.PaymentStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) AddItem(ctx context.Context, userID int32, req *cartDto.CartItemRequest) (*cartDto.CartItemResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cartDto.CartItemResponse), args.Error(1)
}

func (m *MockCartService) UpdateQuantity(ctx context.Context, userID int32, req *cartDto.CartItemRequest) (*cartDto.CartItemResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cartDto.CartItemResponse), args.Error(1)
}

func (m *MockCartService) RemoveItem(ctx context.Context, userID, productID int32) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockCartService) RemoveAllItems(ctx context.Context, userID int32) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCartService) GetItems(ctx context.Context, userID int32) ([]*cartDto.CartItemResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*cartDto.CartItemResponse), args.Error(1)
}

type testSuite struct {
	orderRepo    *MockOrderRepository
	cartService  *MockCartService
	orderService interfaces.OrderService
	ctx          context.Context
}

func setupTest(t *testing.T) *testSuite {
	orderRepo := new(MockOrderRepository)
	cartService := new(MockCartService)
	orderService := NewOrderService(orderRepo, cartService)

	return &testSuite{
		orderRepo:    orderRepo,
		cartService:  cartService,
		orderService: orderService,
		ctx:          context.Background(),
	}
}

func TestOrderService(t *testing.T) {
	t.Run("Create Order - Success", func(t *testing.T) {
		// Setup
		ts := setupTest(t)

		userID := int32(1)
		cartItems := []*cartDto.CartItemResponse{
			{
				ID:        1,
				ProductID: 1,
				Quantity:  2,
				Price:     10.99,
			},
			{
				ID:        2,
				ProductID: 2,
				Quantity:  1,
				Price:     20.99,
			},
		}

		req := &dto.CreateOrderRequest{
			ShippingAddress: "123 Test St",
			ShippingCity:    "Test City",
			ShippingCountry: "Test Country",
			ShippingZip:     "12345",
		}

		// Setup expectations
		ts.cartService.On("GetItems", ts.ctx, userID).Return(cartItems, nil)

		// Calculate expected total
		expectedTotal := 10.99*2 + 20.99

		// Mock order creation
		ts.orderRepo.On("Create", ts.ctx, mock.MatchedBy(func(order *entities.Order) bool {
			return order.UserID == userID &&
				order.TotalAmount == expectedTotal &&
				order.Status == constants.OrderStatusPending &&
				order.PaymentStatus == constants.PaymentStatusPending
		})).Return(&entities.Order{
			ID:              1,
			UserID:          userID,
			Status:          constants.OrderStatusPending,
			PaymentStatus:   constants.PaymentStatusPending,
			TotalAmount:     expectedTotal,
			ShippingAddress: req.ShippingAddress,
			ShippingCity:    req.ShippingCity,
			ShippingCountry: req.ShippingCountry,
			ShippingZip:     req.ShippingZip,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}, nil)

		// Mock order items creation
		ts.orderRepo.On("CreateOrderItems", ts.ctx, int32(1), mock.MatchedBy(func(items []*entities.OrderItem) bool {
			return len(items) == 2 &&
				items[0].ProductID == cartItems[0].ProductID &&
				items[0].Quantity == cartItems[0].Quantity &&
				items[0].Price == cartItems[0].Price &&
				items[1].ProductID == cartItems[1].ProductID &&
				items[1].Quantity == cartItems[1].Quantity &&
				items[1].Price == cartItems[1].Price
		})).Return(nil)

		// Mock cart cleanup
		ts.cartService.On("RemoveAllItems", ts.ctx, userID).Return(nil)

		// Execute test
		order, err := ts.orderService.CreateOrder(ts.ctx, userID, req)
		require.NoError(t, err)
		require.NotNil(t, order)
		require.Equal(t, int32(1), order.ID)
		require.Equal(t, expectedTotal, order.TotalAmount)
		require.Equal(t, constants.OrderStatusPending.String(), order.Status)
		require.Equal(t, constants.PaymentStatusPending.String(), order.PaymentStatus)

		// Verify all expectations
		ts.cartService.AssertExpectations(t)
		ts.orderRepo.AssertExpectations(t)
	})

	t.Run("Create Order - Empty Cart", func(t *testing.T) {
		// Setup new test suite instance
		ts := setupTest(t)

		userID := int32(1)
		req := &dto.CreateOrderRequest{
			ShippingAddress: "123 Test St",
			ShippingCity:    "Test City",
			ShippingCountry: "Test Country",
			ShippingZip:     "12345",
		}

		// Setup expectations for empty cart
		ts.cartService.On("GetItems", ts.ctx, userID).Return([]*cartDto.CartItemResponse{}, nil)

		// Execute test
		order, err := ts.orderService.CreateOrder(ts.ctx, userID, req)
		require.Error(t, err)
		require.Equal(t, errorx.ErrCartEmpty, err)
		require.Nil(t, order)

		ts.cartService.AssertExpectations(t)
	})

	t.Run("Create Order - Failed Cart Cleanup", func(t *testing.T) {
		// Setup new test suite instance
		ts := setupTest(t)

		userID := int32(1)
		cartItems := []*cartDto.CartItemResponse{
			{
				ID:        1,
				ProductID: 1,
				Quantity:  2,
				Price:     10.99,
			},
		}

		req := &dto.CreateOrderRequest{
			ShippingAddress: "123 Test St",
			ShippingCity:    "Test City",
			ShippingCountry: "Test Country",
			ShippingZip:     "12345",
		}

		// Setup expectations
		ts.cartService.On("GetItems", ts.ctx, userID).Return(cartItems, nil)

		// Mock order creation
		ts.orderRepo.On("Create", ts.ctx, mock.Anything).Return(&entities.Order{
			ID:              1,
			UserID:          userID,
			Status:          constants.OrderStatusPending,
			PaymentStatus:   constants.PaymentStatusPending,
			TotalAmount:     21.98, // 10.99 * 2
			ShippingAddress: req.ShippingAddress,
			ShippingCity:    req.ShippingCity,
			ShippingCountry: req.ShippingCountry,
			ShippingZip:     req.ShippingZip,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}, nil)

		// Mock order items creation
		ts.orderRepo.On("CreateOrderItems", ts.ctx, int32(1), mock.Anything).Return(nil)

		// Mock failed cart cleanup
		ts.cartService.On("RemoveAllItems", ts.ctx, userID).Return(errorx.ErrCannotCreateOrder)

		// Execute test
		order, err := ts.orderService.CreateOrder(ts.ctx, userID, req)
		require.NoError(t, err) // Order should still be created even if cart cleanup fails
		require.NotNil(t, order)
		require.Equal(t, int32(1), order.ID)

		// Verify expectations
		ts.cartService.AssertExpectations(t)
		ts.orderRepo.AssertExpectations(t)
	})
}
