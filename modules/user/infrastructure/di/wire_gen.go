// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"mallbots/modules/user/application/services"
	"mallbots/modules/user/infrastructure/repositories"
	"mallbots/modules/user/infrastructure/rest"
)

// Injectors from wire.go:

func InitializeUserHandler(db *pgxpool.Pool) (*rest.UserHandler, error) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := rest.NewUserHandler(userService)
	return userHandler, nil
}

// wire.go:

var UserSet = wire.NewSet(repositories.NewUserRepository, services.NewUserService, rest.NewUserHandler)
