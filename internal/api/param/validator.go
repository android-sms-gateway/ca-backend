package param

import "github.com/gofiber/fiber/v2"

func NewValidator(name string, validator func(string) error) fiber.Handler {
	if name == "" {
		panic("name cannot be empty")
	}
	if validator == nil {
		panic("validator cannot be nil")
	}

	return func(c *fiber.Ctx) error {
		id := c.Params(name)
		if err := validator(id); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.Next()
	}
}
