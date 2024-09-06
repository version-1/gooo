package main

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/version-1/gooo/pkg/app"
	"github.com/version-1/gooo/pkg/config"
	"github.com/version-1/gooo/pkg/controller"
	"github.com/version-1/gooo/pkg/logger"
)

func main() {
	users := controller.GroupHandler{
		Path: "/users",
		Handlers: []controller.Handler{
			{
				Path:   "/",
				Method: http.MethodGet,
			},
			{
				Path:   "/",
				Method: http.MethodPost,
			},
			{
				Path:   "/:id",
				Method: http.MethodPatch,
			},
			{
				Path:   "/:id",
				Method: http.MethodGet,
			},
			{
				Path:   "/:id",
				Method: http.MethodDelete,
			},
		},
	}.List()

	apiRoot := controller.GroupHandler{
		Path: "/api/v1",
	}
	apiRoot.Add(users...)

	cfg := &config.App{
		Logger: logger.DefaultLogger,
	}

	s := app.Server{
		Addr:   ":8080",
		Config: cfg,
	}
	s.RegisterHandlers(apiRoot.List()...)
	app.WithDefaultMiddlewares(&s)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprint(w, logger.DefaultLogger.SInfof("Path\t|\tMethod"))
	s.WalkThrough(func(h controller.Handler) {
		fmt.Fprint(w, logger.DefaultLogger.SInfof("%s\t|\t%s\t", h.Path, h.Method))
	})
	w.Flush()

	s.Run()
}
