package main

import (
	"context"
	"log"
	"net/http"

	"github.com/version-1/gooo/pkg/core/app"
	"github.com/version-1/gooo/pkg/core/request"
	"github.com/version-1/gooo/pkg/core/response"
	"github.com/version-1/gooo/pkg/core/route"
	"github.com/version-1/gooo/pkg/toolkit/logger"
	"github.com/version-1/gooo/pkg/toolkit/middleware"
)

func main() {
	cfg := &app.Config{}
	cfg.SetLogger(logger.DefaultLogger)

	server := &app.App{
		Addr:   ":8080",
		Config: cfg,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			cfg.Logger().Errorf("Error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
		},
	}

	users := route.GroupHandler{
		Path: "/users",
		Handlers: []route.HandlerInterface{
			route.JSON[request.Void, map[string]string]().Get("", func(res *response.Response[map[string]string], req *request.Request[request.Void]) {
				res.Render(map[string]string{"message": "ok"})
			}),
			route.JSON[any, any]().Post("", func(res *response.Response[any], req *request.Request[any]) {
				res.Render(map[string]string{"message": "ok"})
			}),
			route.JSON[request.Void, any]().Get(":id", func(res *response.Response[any], req *request.Request[request.Void]) {
				res.Render(map[string]string{"message": "ok"})
			}),
			route.JSON[request.Void, any]().Patch(":id", func(res *response.Response[any], req *request.Request[request.Void]) {
				res.Render(map[string]string{"message": "ok"})
			}),
			route.JSON[request.Void, any]().Delete(":id", func(res *response.Response[any], req *request.Request[request.Void]) {
				res.Render(map[string]string{"message": "ok"})
			}),
		},
	}
	apiv1 := route.GroupHandler{
		Path: "/api/v1",
	}
	apiv1.Add(users.Children()...)
	app.WithDefaultMiddlewares(server, apiv1.Children()...)
	route.Walk(apiv1.Children(), func(h middleware.Handler) {
		server.Logger().Infof("%s", h.String())
	})

	ctx := context.Background()
	if err := server.Run(ctx); err != nil {
		log.Fatalf("failed to run app: %s", err)
	}
}
