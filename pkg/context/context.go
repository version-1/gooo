package context

import (
	"context"

	"github.com/version-1/gooo/pkg/config"
)

const (
	APP_CONFIG_KEY  = "gooo:request:app_config"
	USER_CONFIG_KEY = "gooo:request:user_config"
)

func Get[T any](ctx context.Context, key string) T {
	return ctx.Value(key).(T)
}

func With[T any](ctx context.Context, key string, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func WithAppConfig(ctx context.Context, cfg *config.App) context.Context {
	return With(ctx, APP_CONFIG_KEY, cfg)
}

func AppConfig(ctx context.Context) *config.App {
	return Get[*config.App](ctx, APP_CONFIG_KEY)
}

func WithUserConfig[T any](ctx context.Context, u T) context.Context {
	return With(ctx, USER_CONFIG_KEY, u)
}

func UserConfig[T any](ctx context.Context) T {
	return Get[T](ctx, USER_CONFIG_KEY)
}
