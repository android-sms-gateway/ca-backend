package logger

import (
	"os"

	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	isDebug := os.Getenv("DEBUG") != ""

	logConfig := zap.NewProductionConfig()
	if isDebug {
		logConfig = zap.NewDevelopmentConfig()
	}

	l, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return l, nil
}
