package config

import (
	"time"

	"github.com/android-sms-gateway/ca/pkg/core/config"
)

type HttpConfig struct {
	Address     string   `envconfig:"HTTP__ADDRESS"`
	ProxyHeader string   `envconfig:"HTTP__PROXY_HEADER"`
	Proxies     []string `envconfig:"HTTP__PROXIES"`
}

type APIConfig struct {
	CORSAllowOrigins string `envconfig:"API__CORS_ALLOW_ORIGINS"`
}

type StorageConfig struct {
	URL string `envconfig:"STORAGE__URL"`
}

type CSR struct {
	TTL time.Duration `envconfig:"CSR__TTL"`
}

type Config struct {
	Http    HttpConfig
	API     APIConfig
	Storage StorageConfig
	CSR     CSR
}

var instance = Config{
	Http: HttpConfig{
		Address: "127.0.0.1:3000",
	},
	API: APIConfig{
		CORSAllowOrigins: "",
	},
	Storage: StorageConfig{
		URL: "redis://localhost:6379/0",
	},
	CSR: CSR{
		TTL: 24 * time.Hour,
	},
}

func New() (Config, error) {
	return instance, config.Load(&instance)
}
