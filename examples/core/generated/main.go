package main

import (
	"context"
	"log"
	"net/http"

	"github.com/version-1/gooo/pkg/core/app"
	"github.com/version-1/gooo/pkg/core/route"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

func main() {
	cfg := &app.Config{}
	cfg.SetLogger(logger.DefaultLogger)

	server := &app.App{
		Addr:   ":8080",
		Config: cfg,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			cfg.Logger().Errorf("Error: %+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		},
	}

	RegisterRoutes(server)
	ctx := context.Background()
	if err := server.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %s", err)
	}
}

func RegisterRoutes(srv *app.App) {
	routes := route.GroupHandler{
		Path:     "/users",
		Handlers: []route.HandlerInterface{
			// ここにルーティングが入ります
		},
	}
	app.WithDefaultMiddlewares(srv, routes.Children()...)
}
