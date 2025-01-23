package http

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func preHandleError(err error, logger *zap.Logger) int {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if code >= 500 {
		logger.Error("An error occurred", zap.Error(err))
	}

	return code
}

type JSONErrorResponse struct {
	Error struct {
		Message string `json:"message"` // Message
		Code    int    `json:"code"`    // Code
	}
}

func NewViewsErrorHandler(logger *zap.Logger, template string, layouts ...string) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := preHandleError(err, logger)

		return c.Status(code).Render(template, fiber.Map{"error": err.Error(), "code": code}, layouts...)
	}
}

func NewJSONErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := preHandleError(err, logger)

		return c.Status(code).JSON(fiber.Map{"error": fiber.Map{"message": err.Error(), "code": code}})
	}
}
