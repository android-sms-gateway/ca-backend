package http

import "github.com/gofiber/fiber/v2"

type Config struct {
	Address     string
	ProxyHeader string
	Proxies     []string
}

type Options struct {
	getOnly      bool
	views        fiber.Views
	errorHandler fiber.ErrorHandler
}

func (o *Options) WithGetOnly() *Options {
	o.getOnly = true
	return o
}

func (o *Options) WithViews(views fiber.Views) *Options {
	o.views = views
	return o
}

func (o *Options) WithErrorHandler(handler fiber.ErrorHandler) *Options {
	o.errorHandler = handler
	return o
}
