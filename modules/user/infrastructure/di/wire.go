//go:build wireinject

package di

import (
	"mallbots/modules/user/application/services"
	"mallbots/modules/user/infrastructure/repositories"
	"mallbots/modules/user/infrastructure/rest"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var UserSet = wire.NewSet(
	repositories.NewUserRepository,
	services.NewUserService,
	rest.NewUserHandler,
)

func InitializeUserHandler(db *pgxpool.Pool) (*rest.UserHandler, error) {
	wire.Build(UserSet)
	return &rest.UserHandler{}, nil
}
