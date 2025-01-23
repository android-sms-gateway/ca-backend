package http

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func New(config Config, option Options, logger *zap.Logger) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		EnableIPValidation:      true,
		EnableTrustedProxyCheck: true,
		ErrorHandler:            option.errorHandler,
		GETOnly:                 option.getOnly,
		ProxyHeader:             config.ProxyHeader,
		TrustedProxies:          config.Proxies,
		Views:                   option.views,
	})
	app.Use(fiberzap.New(fiberzap.Config{
		SkipBody: func(c *fiber.Ctx) bool {
			return c.Response().StatusCode() < 400
		},
		Logger: logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "ua", "body", "error"},
	}))
	app.Use(recover.New())

	return app, nil
}
