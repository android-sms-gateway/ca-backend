package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// New returns a middleware that limits the number of requests from a client
// in a given duration.
//
// max is the maximum number of requests allowed in the duration.
// duration is the duration of the limit.
//
// If max is less than or equal to 0 or duration is less than or equal to 0,
// a panic occurs.
func New(max int, duration time.Duration) fiber.Handler {
	if max <= 0 {
		panic("max must be greater than 0")
	}
	if duration <= 0 {
		panic("duration must be positive")
	}

	return limiter.New(limiter.Config{
		Max:                max,
		SkipFailedRequests: true,
		Expiration:         duration,
	})
}
