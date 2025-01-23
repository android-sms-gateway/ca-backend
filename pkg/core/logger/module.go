package logger

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"logger",
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				_ = logger.Sync()
				return nil
			},
		})
	}),
)
