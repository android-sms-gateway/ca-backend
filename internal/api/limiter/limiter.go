package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func New(max int) fiber.Handler {
	if max <= 0 {
		panic("max must be greater than 0")
	}

	return limiter.New(limiter.Config{
		Max:                max,
		SkipFailedRequests: true,
		Expiration:         time.Second,
	})
}
