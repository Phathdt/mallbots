package interfaces

import (
	"context"
	"mallbots/modules/product/domain/entities"

	"github.com/phathdt/service-context/core"
)

type ProductRepository interface {
	GetProducts(ctx context.Context, filter *ProductFilter, paging *core.Paging) ([]*entities.Product, error)
	GetProduct(ctx context.Context, id int32) (*entities.Product, error)
}

type ProductFilter struct {
	Search   string
	MinPrice *float64
	MaxPrice *float64
	Category *int32
	SortBy   string
}
