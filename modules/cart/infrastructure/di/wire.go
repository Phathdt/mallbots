//go:build wireinject

package di

import (
	"mallbots/modules/cart/application/services"
	"mallbots/modules/cart/infrastructure/repositories"
	"mallbots/modules/cart/infrastructure/rest"
	productService "mallbots/modules/product/application/services"
	productRepo "mallbots/modules/product/infrastructure/repositories"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var CartSet = wire.NewSet(
	productRepo.NewProductRepository,
	productService.NewProductService,
	repositories.NewCartRepository,
	services.NewCartService,
	rest.NewCartHandler,
)

func InitializeCartHandler(db *pgxpool.Pool) (*rest.CartHandler, error) {
	wire.Build(CartSet)
	return &rest.CartHandler{}, nil
}
