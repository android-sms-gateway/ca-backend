package api

import (
	"github.com/android-sms-gateway/ca/pkg/core/http"
	"github.com/android-sms-gateway/ca/pkg/core/http/jsonify"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	fx.Invoke(func(app *fiber.App, config Config) {
		api := app.Group("/api/v1")

		apidoc.SwaggerInfo.Version = version.AppVersion
		api.Get("/docs/*", swagger.New(swagger.Config{
			Layout: "BaseLayout",
		}))

		if config.CORSAllowOrigins != "" {
			api.Use(cors.New(cors.Config{
				AllowOrigins:     config.CORSAllowOrigins,
				AllowCredentials: true,
				MaxAge:           86400,
			}))
		}

		api.Use(jsonify.New())

		api.Use(func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusNotFound)
		})
	}),
)
