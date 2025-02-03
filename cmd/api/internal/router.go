package cmd

import (
	"log"
	"log/slog"
	"mallbots/modules/product/infrastructure/di"
	"mallbots/plugins/pgxc"
	"mallbots/shared/common"
	"mallbots/shared/config"
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

	productHandler, err := di.InitializeProductHandler(dbPool)
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

	_ = app.Listen(":4000")
}

func ping() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(&fiber.Map{
			"msg": "pong",
		})
	}
}
