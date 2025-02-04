package services

import (
	"context"
	"mallbots/modules/cart/application/dto"
	"mallbots/modules/cart/domain/entities"
	"mallbots/modules/cart/domain/interfaces"
	productInterfaces "mallbots/modules/product/domain/interfaces"
	"time"
)

type cartService struct {
	cartRepo       interfaces.CartRepository
	productService productInterfaces.ProductService
}

func NewCartService(
	cartRepo interfaces.CartRepository,
	productService productInterfaces.ProductService,
) interfaces.CartService {
	return &cartService{
		cartRepo:       cartRepo,
		productService: productService,
	}
}

func (s *cartService) AddItem(ctx context.Context, userID int32, req *dto.CartItemRequest) (*dto.CartItemResponse, error) {
	// Get product to validate and get current price
	product, err := s.productService.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.GetByUserAndProduct(ctx, userID, req.ProductID)
	if err == nil && existingItem != nil {
		// Update quantity if item exists
		existingItem.Quantity += req.Quantity
		existingItem.UpdatedAt = time.Now()

		if err := s.cartRepo.Update(ctx, existingItem); err != nil {
			return nil, err
		}

		return &dto.CartItemResponse{
			ID:        existingItem.ID,
			ProductID: existingItem.ProductID,
			Quantity:  existingItem.Quantity,
			Price:     existingItem.Price,
		}, nil
	}

	// Create new cart item
	cartItem := &entities.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     product.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newItem, err := s.cartRepo.Create(ctx, cartItem)
	if err != nil {
		return nil, err
	}

	return &dto.CartItemResponse{
		ID:        newItem.ID,
		ProductID: newItem.ProductID,
		Quantity:  newItem.Quantity,
		Price:     newItem.Price,
	}, nil
}

func (s *cartService) UpdateQuantity(ctx context.Context, userID int32, req *dto.CartItemRequest) (*dto.CartItemResponse, error) {
	item, err := s.cartRepo.GetByUserAndProduct(ctx, userID, req.ProductID)
	if err != nil {
		return nil, err
	}

	item.Quantity = req.Quantity
	item.UpdatedAt = time.Now()

	if err := s.cartRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return &dto.CartItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
	}, nil
}

func (s *cartService) RemoveItem(ctx context.Context, userID, productID int32) error {
	return s.cartRepo.Delete(ctx, userID, productID)
}

func (s *cartService) GetItems(ctx context.Context, userID int32) ([]*dto.CartItemResponse, error) {
	items, err := s.cartRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var response []*dto.CartItemResponse
	for _, item := range items {
		response = append(response, &dto.CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return response, nil
}

func (s *cartService) RemoveAllItems(ctx context.Context, userID int32) error {
	return s.cartRepo.DeleteAllByUser(ctx, userID)
}
