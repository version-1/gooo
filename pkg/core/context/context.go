package context

import (
	"context"
)

const (
	APP_CONFIG_KEY = "gooo:request:app_config"
)

func Get[T any](ctx context.Context, key string) T {
	return ctx.Value(key).(T)
}

func With[T any](ctx context.Context, key string, value T) context.Context {
	return context.WithValue(ctx, key, value)
}
