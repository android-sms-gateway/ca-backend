package api

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/core/http/jsonify"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"

	apidoc "github.com/android-sms-gateway/ca/api"
	"github.com/android-sms-gateway/ca/internal/version"
)

var Module = fx.Module(
	"api",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("api")
	}),
	fx.Provide(http.NewJSONErrorHandler),
	fx.Provide(func(log *zap.Logger) http.Options {
		return *(&http.Options{}).WithErrorHandler(http.NewJSONErrorHandler(log))
	}),
	fx.Provide(newCSR, fx.Private),
	fx.Invoke(func(app *fiber.App, csr *csrHandler, config Config) {
		apidoc.SwaggerInfo.Version = version.AppVersion
		app.Use("/docs/*",
			etag.New(etag.Config{
				Weak: true,
			}),
			swagger.New(swagger.Config{
				Layout: "BaseLayout",
			}),
		)

		metrics := fiberprometheus.New("")
		metrics.RegisterAt(app, "/metrics")
		app.Use(metrics.Middleware)

		api := app.Group("/api/v1")

		if config.CORSAllowOrigins != "" {
			api.Use(cors.New(cors.Config{
				AllowOrigins:     config.CORSAllowOrigins,
				AllowCredentials: true,
				MaxAge:           86400,
			}))
		}

		api.Use(jsonify.New())

		csr.Register(api.Group("csr"))

		app.Use(func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusNotFound)
		})
	}),
)
