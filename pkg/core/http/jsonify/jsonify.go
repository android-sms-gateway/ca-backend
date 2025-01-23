package jsonify

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return err
		}

		contentType := string(c.Response().Header.ContentType())
		if strings.Contains(contentType, "application/json") {
			return nil
		}

		body := c.Response().Body()

		if c.Response().StatusCode() < 400 {
			return c.JSON(body)
		}

		return fiber.NewError(c.Response().StatusCode(), string(body))
	}
}
