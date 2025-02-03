package services

import (
	"context"
	"mallbots/modules/product/application/dto"
	"mallbots/modules/product/domain/interfaces"

	"github.com/phathdt/service-context/core"
)

type ProductService struct {
	repo interfaces.ProductRepository
}

func NewProductService(repo interfaces.ProductRepository) interfaces.ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts(ctx context.Context, req *dto.ProductListRequest, paging *core.Paging) ([]*dto.ProductResponse, error) {
	filter := interfaces.ProductFilter{
		Search:   req.Search,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		Category: req.Category,
		SortBy:   req.SortBy,
	}

	products, err := s.repo.GetProducts(ctx, &filter, paging)
	if err != nil {
		return nil, err
	}

	var response []*dto.ProductResponse
	for _, p := range products {
		response = append(response, &dto.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CategoryID:  p.CategoryID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return response, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id int32) (*dto.ProductResponse, error) {
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}
