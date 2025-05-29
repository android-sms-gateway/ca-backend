package internal

import (
	"github.com/android-sms-gateway/ca/internal/api"
	"github.com/android-sms-gateway/ca/internal/config"
	"github.com/android-sms-gateway/ca/internal/csr"
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/core/logger"
	"github.com/android-sms-gateway/core/redis"
	"github.com/android-sms-gateway/core/validator"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	fx.New(
		logger.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
		http.Module,
		validator.Module,
		redis.Module,

		config.Module,
		api.Module,
		csr.Module,
	).
		Run()
}
