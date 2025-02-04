//go:build wireinject

package di

import (
	"mallbots/modules/cart/application/services"
	"mallbots/modules/cart/infrastructure/repositories"
	"mallbots/modules/cart/infrastructure/rest"
	productRepo "mallbots/modules/product/infrastructure/repositories"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var CartSet = wire.NewSet(
	repositories.NewCartItemRepository,
	productRepo.NewProductRepository,
	services.NewCartItemService,
	rest.NewCartItemHandler,
)

func InitializeCartHandler(db *pgxpool.Pool) (*rest.CartItemHandler, error) {
	wire.Build(CartSet)
	return &rest.CartItemHandler{}, nil
}
