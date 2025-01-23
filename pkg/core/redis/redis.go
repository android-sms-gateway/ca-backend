package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func New(c Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(c.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	return redis.NewClient(opts), nil
}
