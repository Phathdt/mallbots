//go:build wireinject

package di

import (
	"mallbots/modules/product/application/services"
	"mallbots/modules/product/infrastructure/repositories"
	"mallbots/modules/product/infrastructure/rest"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ProductSet = wire.NewSet(
	repositories.NewProductRepository,
	services.NewProductService,
	rest.NewProductHandler,
)

func InitializeProductHandler(db *pgxpool.Pool) (*rest.ProductHandler, error) {
	wire.Build(ProductSet)
	return &rest.ProductHandler{}, nil
}
