package cmd

import (
	"log"
	"log/slog"
	cartDi "mallbots/modules/cart/infrastructure/di"
	orderDi "mallbots/modules/order/infrastructure/di"
	productDi "mallbots/modules/product/infrastructure/di"
	userDi "mallbots/modules/user/infrastructure/di"
	"mallbots/plugins/pgxc"
	"mallbots/plugins/tokenprovider"
	"mallbots/shared/common"
	"mallbots/shared/config"
	middleware2 "mallbots/shared/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	sctx "github.com/phathdt/service-context"
	"github.com/phathdt/service-context/component/fiberc/middleware"
	slogfiber "github.com/samber/slog-fiber"
)

func StartRouter(sc sctx.ServiceContext, cfg *config.Config) {
	dbPool := sc.MustGet(common.KeyPgx).(pgxc.PgxComp).GetConn()

	tokenProvider := sc.MustGet(common.KeyJwt).(tokenprovider.Provider)

	productHandler, err := productDi.InitializeProductHandler(dbPool)
	if err != nil {
		log.Fatal(err)
	}

	userHandler, err := userDi.InitializeUserHandler(dbPool, tokenProvider)
	if err != nil {
		log.Fatal(err)
	}

	cartHandler, err := cartDi.InitializeCartHandler(dbPool)
	if err != nil {
		log.Fatal(err)
	}

	orderHandler, err := orderDi.InitializeOrderHandler(dbPool)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{BodyLimit: 100 * 1024 * 1024})

	app.Use(slogfiber.New(slog.New(slog.NewTextHandler(os.Stdout, nil))))
	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(middleware.Recover(sc))

	app.Get("/", ping())

	// Setup routes
	app.Get("/v1/products", productHandler.GetProducts)
	app.Get("/v1/products/:id", productHandler.GetProduct)

	// User routes
	app.Post("/v1/auth/register", userHandler.Register)
	app.Post("/v1/auth/login", userHandler.Login)

	// Protected routes
	app.Use(middleware2.RequiredAuth(sc))

	app.Get("/v1/users/me", userHandler.GetProfile)

	// Cart routes
	app.Post("/v1/cart/items", cartHandler.AddItem)
	app.Put("/v1/cart/items", cartHandler.UpdateQuantity)
	app.Delete("/v1/cart/items/:productId", cartHandler.RemoveItem)
	app.Get("/v1/cart/items", cartHandler.GetItems)

	// Order routes
	app.Post("/v1/orders", orderHandler.CreateOrder)
	app.Get("/v1/orders", orderHandler.GetUserOrders)
	app.Get("/v1/orders/:id", orderHandler.GetOrder)

	_ = app.Listen(":4000")
}

func ping() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(&fiber.Map{
			"msg": "pong",
		})
	}
}
