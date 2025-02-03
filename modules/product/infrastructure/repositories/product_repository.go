package repositories

import (
	"context"
	"mallbots/modules/product/domain/entities"
	"mallbots/modules/product/domain/interfaces"
	"mallbots/modules/product/infrastructure/query/gen"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phathdt/service-context/core"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) interfaces.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetProducts(ctx context.Context, filter *interfaces.ProductFilter, paging *core.Paging) ([]*entities.Product, error) {
	queries := gen.New(r.db)

	offset := (paging.Page - 1) * paging.Limit

	// Default values for nullable fields
	categoryID := int32(0)
	if filter.Category != nil {
		categoryID = *filter.Category
	}

	minPrice := float64(0)
	if filter.MinPrice != nil {
		minPrice = *filter.MinPrice
	}

	maxPrice := float64(0)
	if filter.MaxPrice != nil {
		maxPrice = *filter.MaxPrice
	}

	// Get total count for pagination
	total, err := queries.CountProducts(ctx, gen.CountProductsParams{
		Btrim:   filter.Search,
		Column2: categoryID,
		Column3: minPrice,
		Column4: maxPrice,
	})
	if err != nil {
		return nil, err
	}
	paging.Total = int64(total)

	// Get products with pagination
	products, err := queries.GetProducts(ctx, gen.GetProductsParams{
		Btrim:   filter.Search,
		Column2: categoryID,
		Column3: minPrice,
		Column4: maxPrice,
		Column5: filter.SortBy,
		Limit:   int32(paging.Limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Product, len(products))
	for i, p := range products {
		result[i] = &entities.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CategoryID:  p.CategoryID,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	return result, nil
}

func (r *productRepository) GetProduct(ctx context.Context, id int32) (*entities.Product, error) {
	queries := gen.New(r.db)

	product, err := queries.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entities.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}
