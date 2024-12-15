package context

import (
	"context"
	"errors"
	"fmt"
)

const (
	APP_CONFIG_KEY = "gooo:request:app_config"
)

func Get[T any](ctx context.Context, key string) T {
	v, ok := ctx.Value(key).(T)
	if !ok {
		err := errors.New(fmt.Sprintf("context value not found: %s", key))
		panic(err)
	}

	return v
}

func With[T any](ctx context.Context, key string, value T) context.Context {
	return context.WithValue(ctx, key, value)
}
