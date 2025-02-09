package main

// This is a generated file. DO NOT EDIT manually.
import (
	"context"
	"log"
	"net/http"

	"github.com/version-1/gooo/examples/core/internal/schema"
	"github.com/version-1/gooo/pkg/core/app"
	"github.com/version-1/gooo/pkg/core/request"
	"github.com/version-1/gooo/pkg/core/response"
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
		Path: "/users",
		Handlers: []route.HandlerInterface{
			route.JSON[request.Void, schema.User]().Get("/users", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[schema.MutateUser, schema.User]().Post("/users", func(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
				// do something
			}),
			route.JSON[request.Void, schema.User]().Delete("/users/{id}", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[request.Void, schema.User]().Get("/users/{id}", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[schema.MutateUser, schema.User]().Patch("/users/{id}", func(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
				// do something
			}),
			route.JSON[request.Void, schema.Post]().Get("/posts", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[schema.MutatePost, schema.Post]().Post("/posts", func(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
				// do something
			}),
			route.JSON[request.Void, schema.Post]().Get("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[schema.MutatePost, schema.Post]().Patch("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
				// do something
			}),
			route.JSON[request.Void, schema.Post]().Delete("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
				// do something
			}),
		},
	}
	app.WithDefaultMiddlewares(srv, routes.Children()...)
}
