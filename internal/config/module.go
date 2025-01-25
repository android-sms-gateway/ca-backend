package config

import (
	"github.com/android-sms-gateway/ca/internal/api"
	"github.com/android-sms-gateway/ca/internal/csr"
	"github.com/android-sms-gateway/ca/pkg/core/http"
	"github.com/android-sms-gateway/ca/pkg/core/redis"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(New),
	fx.Provide(func(c Config) http.Config {
		return http.Config{
			Address:     c.Http.Address,
			ProxyHeader: c.Http.ProxyHeader,
			Proxies:     c.Http.Proxies,
		}
	}),
	fx.Provide(func(c Config) redis.Config {
		return redis.Config{
			URL: c.Storage.URL,
		}
	}),
	fx.Provide(func(c Config) api.Config {
		return api.Config{
			CORSAllowOrigins: c.API.CORSAllowOrigins,
		}
	}),
	fx.Provide(func(c Config) csr.Config {
		return csr.Config{
			TTL: c.CSR.TTL,
		}
	}),
)
