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
	fx.Provide(func(config Config, csrs *repository, log *zap.Logger) (*Service, error) {
		caCert, caKey, err := LoadCA(config.CACert, config.CAKey)
		if err != nil {
			return nil, err
		}

		return NewService(csrs, caCert, caKey, log), nil
	}),
	fx.Invoke(func(lc fx.Lifecycle, s *Service) {
		lc.Append(fx.Hook{
			OnStop: s.Stop,
		})
	}),
)
