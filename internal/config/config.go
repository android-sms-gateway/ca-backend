package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-core-fx/config"
)

type httpConfig struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`
}

type apiConfig struct {
	CORSAllowOrigins string `koanf:"cors_allow_origins"`
}

type storageConfig struct {
	URL string `koanf:"url"`
}

type csrConfig struct {
	TTL        time.Duration `koanf:"ttl"`
	CACertPath string        `koanf:"ca_cert_path"`
	CAKeyPath  string        `koanf:"ca_key_path"`
}

type Config struct {
	HTTP    httpConfig    `koanf:"http"`
	API     apiConfig     `koanf:"api"`
	Storage storageConfig `koanf:"storage"`
	CSR     csrConfig     `koanf:"csr"`
}

func Default() Config {
	//nolint:mnd //default values
	return Config{
		HTTP: httpConfig{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "",
			Proxies:     []string{},
		},
		API: apiConfig{
			CORSAllowOrigins: "",
		},
		Storage: storageConfig{
			URL: "redis://localhost:6379/0",
		},
		CSR: csrConfig{
			TTL:        24 * time.Hour,
			CACertPath: "",
			CAKeyPath:  "",
		},
	}
}

func New() (Config, error) {
	cfg := Default()

	options := []config.Option{}
	if yamlPath := os.Getenv("CONFIG_PATH"); yamlPath != "" {
		options = append(options, config.WithLocalYAML(yamlPath))
	}

	if err := config.Load(&cfg, options...); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
