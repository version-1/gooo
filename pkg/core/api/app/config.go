package app

import (
	"github.com/gooolib/logger"
)

type Config struct {
	logger logger.Logger
}

func (c *Config) SetLogger(l logger.Logger) {
	c.logger = l
}

func (c Config) Logger() logger.Logger {
	if c.logger == nil {
		return logger.DefaultLogger
	}

	return c.logger
}
