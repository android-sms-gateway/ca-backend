package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Base struct {
	Validator *validator.Validate
	Logger    *zap.Logger
}

func (c *Base) BodyParserValidator(ctx *fiber.Ctx, out interface{}) error {
	if err := ctx.BodyParser(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to parse request: %s", err.Error()))
	}

	if err := c.Validator.Struct(out); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to validate request: %s", err.Error()))
	}

	if v, ok := out.(Validatable); ok {
		if err := v.Validate(); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to validate request: %s", err.Error()))
		}
	}

	return nil
}
