package config

import (
	"github.com/version-1/gooo/pkg/logger"
)

type App struct {
	Logger                  logger.Logger
	DefaultResponseRenderer string
}

type ResponseRenderer string

const (
	JSONAPIRenderer ResponseRenderer = "jsonapi"
	RawRenderer     ResponseRenderer = "raw"
)

func (c App) GetLogger() logger.Logger {
	if c.Logger == nil {
		return logger.DefaultLogger
	}

	return c.Logger
}
