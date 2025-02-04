package services

import (
	"context"
	"errors"
	"mallbots/modules/cart/application/dto"
	"mallbots/modules/cart/domain/entities"
	productDto "mallbots/modules/product/application/dto"
	"testing"
	"time"

	"github.com/phathdt/service-context/core"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock repositories
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) Create(ctx context.Context, item *entities.CartItem) (*entities.CartItem, error) {
	args := m.Called(ctx, item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CartItem), args.Error(1)
}

func (m *MockCartRepository) Update(ctx context.Context, item *entities.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockCartRepository) Delete(ctx context.Context, userID, productID int32) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockCartRepository) DeleteAllByUser(ctx context.Context, userID int32) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCartRepository) GetByUserAndProduct(ctx context.Context, userID, productID int32) (*entities.CartItem, error) {
	args := m.Called(ctx, userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CartItem), args.Error(1)
}

func (m *MockCartRepository) GetByUser(ctx context.Context, userID int32) ([]*entities.CartItem, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.CartItem), args.Error(1)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProduct(ctx context.Context, id int32) (*productDto.ProductResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productDto.ProductResponse), args.Error(1)
}

func (m *MockProductService) GetProducts(ctx context.Context, req *productDto.ProductListRequest, paging *core.Paging) ([]*productDto.ProductResponse, error) {
	args := m.Called(ctx, req, paging)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*productDto.ProductResponse), args.Error(1)
}

func TestCartService(t *testing.T) {
	ctx := context.Background()
	cartRepo := new(MockCartRepository)
	productService := new(MockProductService)
	cartService := NewCartService(cartRepo, productService)

	t.Run("Add Item to Cart", func(t *testing.T) {
		userID := int32(1)
		req := &dto.CartItemRequest{
			ProductID: 1,
			Quantity:  2,
		}

		// Mock product service response
		productService.On("GetProduct", ctx, req.ProductID).Return(&productDto.ProductResponse{
			ID:    1,
			Price: 10.99,
		}, nil)

		// Mock repository calls
		cartRepo.On("GetByUserAndProduct", ctx, userID, req.ProductID).Return(nil, errors.New("not found"))
		cartRepo.On("Create", ctx, mock.AnythingOfType("*entities.CartItem")).Return(&entities.CartItem{
			ID:        1,
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Price:     10.99,
		}, nil)

		// Test add item
		response, err := cartService.AddItem(ctx, userID, req)
		require.NoError(t, err)
		require.Equal(t, req.ProductID, response.ProductID)
		require.Equal(t, req.Quantity, response.Quantity)
	})

	t.Run("Remove All Items from Cart", func(t *testing.T) {
		userID := int32(1)

		// Mock repository call
		cartRepo.On("DeleteAllByUser", ctx, userID).Return(nil)

		// Test remove all items
		err := cartService.RemoveAllItems(ctx, userID)
		require.NoError(t, err)

		cartRepo.AssertCalled(t, "DeleteAllByUser", ctx, userID)
	})

	t.Run("Get Cart Items", func(t *testing.T) {
		userID := int32(1)
		mockItems := []*entities.CartItem{
			{
				ID:        1,
				UserID:    userID,
				ProductID: 1,
				Quantity:  2,
				Price:     10.99,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				UserID:    userID,
				ProductID: 2,
				Quantity:  1,
				Price:     20.99,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Mock repository call
		cartRepo.On("GetByUser", ctx, userID).Return(mockItems, nil)

		// Test get items
		items, err := cartService.GetItems(ctx, userID)
		require.NoError(t, err)
		require.Len(t, items, 2)
		require.Equal(t, mockItems[0].ID, items[0].ID)
		require.Equal(t, mockItems[0].ProductID, items[0].ProductID)
		require.Equal(t, mockItems[0].Quantity, items[0].Quantity)
		require.Equal(t, mockItems[0].Price, items[0].Price)
	})
}
