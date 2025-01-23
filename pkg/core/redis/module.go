package redis

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"redis",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("redis")
	}),
	fx.Provide(New),
)
