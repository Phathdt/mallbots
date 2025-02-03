package interfaces

import (
	"context"
	"mallbots/modules/product/application/dto"

	"github.com/phathdt/service-context/core"
)

type ProductService interface {
	GetProducts(ctx context.Context, req *dto.ProductListRequest, paging *core.Paging) ([]*dto.ProductResponse, error)
	GetProduct(ctx context.Context, id int32) (*dto.ProductResponse, error)
}
