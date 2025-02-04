package services

import (
	"context"
	"errors"
	"mallbots/modules/cart/application/dto"
	"mallbots/modules/cart/domain/entities"
	productDTO "mallbots/modules/product/application/dto"
	"testing"
	"time"

	"github.com/phathdt/service-context/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) Create(ctx context.Context, item *entities.CartItem) (*entities.CartItem, error) {
	args := m.Called(ctx, item)
	if v, ok := args.Get(0).(*entities.CartItem); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCartRepository) Update(ctx context.Context, item *entities.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockCartRepository) Delete(ctx context.Context, userID, productID int32) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockCartRepository) GetByUserAndProduct(ctx context.Context, userID, productID int32) (*entities.CartItem, error) {
	args := m.Called(ctx, userID, productID)
	if v, ok := args.Get(0).(*entities.CartItem); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCartRepository) GetByUser(ctx context.Context, userID int32) ([]*entities.CartItem, error) {
	args := m.Called(ctx, userID)
	if v, ok := args.Get(0).([]*entities.CartItem); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProduct(ctx context.Context, id int32) (*productDTO.ProductResponse, error) {
	args := m.Called(ctx, id)
	if v, ok := args.Get(0).(*productDTO.ProductResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductService) GetProducts(ctx context.Context, req *productDTO.ProductListRequest, paging *core.Paging) ([]*productDTO.ProductResponse, error) {
	args := m.Called(ctx, req, paging)
	if v, ok := args.Get(0).([]*productDTO.ProductResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestCartService_AddItem(t *testing.T) {
	tests := []struct {
		name          string
		userID        int32
		request       *dto.CartItemRequest
		setupMocks    func(*MockCartRepository, *MockProductService)
		expectedItem  *dto.CartItemResponse
		expectedError error
	}{
		{
			name:   "Successfully add new item to cart",
			userID: 1,
			request: &dto.CartItemRequest{
				ProductID: 1,
				Quantity:  2,
			},
			setupMocks: func(mockCartRepo *MockCartRepository, mockProductService *MockProductService) {
				// Setup product service mock
				mockProductService.On("GetProduct", mock.Anything, int32(1)).Return(&productDTO.ProductResponse{
					ID:    1,
					Price: 100.0,
				}, nil)

				// Setup cart repo mocks
				mockCartRepo.On("GetByUserAndProduct", mock.Anything, int32(1), int32(1)).Return(nil, errors.New("not found"))
				mockCartRepo.On("Create", mock.Anything, mock.MatchedBy(func(item *entities.CartItem) bool {
					return item.UserID == 1 && item.ProductID == 1 && item.Quantity == 2 && item.Price == 100.0
				})).Return(&entities.CartItem{
					ID:        1,
					UserID:    1,
					ProductID: 1,
					Quantity:  2,
					Price:     100.0,
				}, nil)
			},
			expectedItem: &dto.CartItemResponse{
				ID:        1,
				ProductID: 1,
				Quantity:  2,
				Price:     100.0,
			},
			expectedError: nil,
		},
		{
			name:   "Update existing item in cart",
			userID: 1,
			request: &dto.CartItemRequest{
				ProductID: 1,
				Quantity:  2,
			},
			setupMocks: func(mockCartRepo *MockCartRepository, mockProductService *MockProductService) {
				existingItem := &entities.CartItem{
					ID:        1,
					UserID:    1,
					ProductID: 1,
					Quantity:  3,
					Price:     100.0,
				}

				// Setup product service mock
				mockProductService.On("GetProduct", mock.Anything, int32(1)).Return(&productDTO.ProductResponse{
					ID:    1,
					Price: 100.0,
				}, nil)

				// Setup cart repo mocks
				mockCartRepo.On("GetByUserAndProduct", mock.Anything, int32(1), int32(1)).Return(existingItem, nil)
				mockCartRepo.On("Update", mock.Anything, mock.MatchedBy(func(item *entities.CartItem) bool {
					return item.ID == 1 && item.Quantity == 5 // 3 + 2
				})).Return(nil)
			},
			expectedItem: &dto.CartItemResponse{
				ID:        1,
				ProductID: 1,
				Quantity:  5,
				Price:     100.0,
			},
			expectedError: nil,
		},
		{
			name:   "Product not found",
			userID: 1,
			request: &dto.CartItemRequest{
				ProductID: 1,
				Quantity:  2,
			},
			setupMocks: func(mockCartRepo *MockCartRepository, mockProductService *MockProductService) {
				mockProductService.On("GetProduct", mock.Anything, int32(1)).Return(nil, errors.New("product not found"))
			},
			expectedItem:  nil,
			expectedError: errors.New("product not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockCartRepo := new(MockCartRepository)
			mockProductService := new(MockProductService)
			tt.setupMocks(mockCartRepo, mockProductService)

			service := NewCartService(mockCartRepo, mockProductService)

			// Execute
			item, err := service.AddItem(context.Background(), tt.userID, tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItem, item)
			}

			// Verify all mocks
			mockCartRepo.AssertExpectations(t)
			mockProductService.AssertExpectations(t)
		})
	}
}

// Helper function to create a CartItem with current time
func createCartItem(id, userID, productID, quantity int32, price float64) *entities.CartItem {
	now := time.Now()
	return &entities.CartItem{
		ID:        id,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
