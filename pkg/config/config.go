package config

import (
	"github.com/version-1/gooo/pkg/logger"
)

type App struct {
	Logger logger.Logger
}

func (c App) GetLogger() logger.Logger {
	if c.Logger == nil {
		return logger.DefaultLogger
	}

	return c.Logger
}
