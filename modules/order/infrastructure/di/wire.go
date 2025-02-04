//go:build wireinject

package di

import (
	cartService "mallbots/modules/cart/application/services"
	cartRepo "mallbots/modules/cart/infrastructure/repositories"
	"mallbots/modules/order/application/services"
	"mallbots/modules/order/infrastructure/repositories"
	"mallbots/modules/order/infrastructure/rest"
	productService "mallbots/modules/product/application/services"
	productRepo "mallbots/modules/product/infrastructure/repositories"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var OrderSet = wire.NewSet(
	productRepo.NewProductRepository,
	productService.NewProductService,
	cartRepo.NewCartRepository,
	cartService.NewCartService,
	repositories.NewOrderRepository,
	services.NewOrderService,
	rest.NewOrderHandler,
)

func InitializeOrderHandler(db *pgxpool.Pool) (*rest.OrderHandler, error) {
	wire.Build(OrderSet)
	return &rest.OrderHandler{}, nil
}
