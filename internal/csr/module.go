package csr

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"csr",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("csr")
	}),
	fx.Provide(func(config Config, redis *redis.Client) *repository {
		return newRepository(redis, config.TTL)
	}, fx.Private),
	fx.Provide(NewService),
	fx.Invoke(func(lc fx.Lifecycle, s *Service) {
		lc.Append(fx.Hook{
			OnStop: s.Stop,
		})
	}),
)
