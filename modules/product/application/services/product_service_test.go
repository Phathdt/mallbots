package services

import (
	"context"
	"mallbots/modules/product/application/dto"
	"mallbots/modules/product/domain/entities"
	"mallbots/modules/product/domain/interfaces"
	"testing"

	"github.com/phathdt/service-context/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) GetProducts(ctx context.Context, filter *interfaces.ProductFilter, paging *core.Paging) ([]*entities.Product, error) {
	args := m.Called(ctx, filter, paging)
	return args.Get(0).([]*entities.Product), args.Error(1)
}

func (m *MockProductRepo) GetProduct(ctx context.Context, id int32) (*entities.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Product), args.Error(1)
}

func TestGetProduct(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepo)
	service := NewProductService(mockRepo)

	expectedProduct := &entities.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      99.99,
		CategoryID: 1,
	}

	mockRepo.On("GetProduct", mock.Anything, int32(1)).Return(expectedProduct, nil)

	// Act
	result, err := service.GetProduct(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProduct.Name, result.Name)
	assert.Equal(t, expectedProduct.Price, result.Price)
	mockRepo.AssertExpectations(t)
}

func TestGetProducts(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepo)
	service := NewProductService(mockRepo)

	req := &dto.ProductListRequest{
		Search:   "test",
		MinPrice: &[]float64{10.0}[0],
		MaxPrice: &[]float64{100.0}[0],
	}

	paging := &core.Paging{
		Page:  1,
		Limit: 10,
	}

	expectedProducts := []*entities.Product{
		{ID: 1, Name: "Test 1", Price: 50.0},
		{ID: 2, Name: "Test 2", Price: 75.0},
	}

	mockRepo.On("GetProducts", mock.Anything, mock.Anything, paging).Return(expectedProducts, nil)

	// Act
	results, err := service.GetProducts(context.Background(), req, paging)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, expectedProducts[0].Name, results[0].Name)
	mockRepo.AssertExpectations(t)
}
