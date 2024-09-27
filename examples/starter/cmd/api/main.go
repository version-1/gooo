package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/version-1/gooo/pkg/app"
	"github.com/version-1/gooo/pkg/config"
	"github.com/version-1/gooo/pkg/controller"
	"github.com/version-1/gooo/pkg/http/request"
	"github.com/version-1/gooo/pkg/http/response"
	"github.com/version-1/gooo/pkg/logger"
	"github.com/version-1/gooo/pkg/presenter/jsonapi"
)

type Dummy struct {
	ID     string    `json:"-"`
	String string    `json:"string"`
	Number int       `json:"number"`
	Flag   bool      `json:"flag"`
	Time   time.Time `json:"time"`
}

func (e Dummy) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	r := jsonapi.Resource{
		ID:         e.ID,
		Type:       "dummy",
		Attributes: jsonapi.NewAttributes(e),
	}

	return r, jsonapi.Resources{}
}

type DummyError struct{}

func (e DummyError) Error() string {
	return "overridden error"
}

func (e DummyError) Code() string {
	return "overridden_error"
}

func (e DummyError) Title() string {
	return "Overrridden Error"
}

func main() {
	ping := controller.Handler{
		Path:   "/ping",
		Method: http.MethodGet,
		Handler: func(w *response.Response, r *request.Request) {
			w.JSON(map[string]string{"message": "pong"})
		},
	}

	testing := controller.GroupHandler{
		Path: "/testing",
		Handlers: []controller.Handler{
			{
				Path:   "/json",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					w.JSON(map[string]string{"message": "ok"})
				},
			},
			{
				Path:   "/render",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					data := Dummy{
						ID:     "1",
						String: "Hello, World!",
						Number: 42,
						Flag:   true,
						Time:   time.Now(),
					}
					if err := w.Render(data); err != nil {
						fmt.Printf("error: %+v\n", err)
						w.InternalServerErrorWith(err)
					}
				},
			},
			{
				Path:   "/render_error",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					if err := w.RenderError(fmt.Errorf("error")); err != nil {
						w.InternalServerErrorWith(err)
					}
				},
			},
			{
				Path:   "/internal_server_error",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					w.InternalServerErrorWith(DummyError{})
				},
			},
			{
				Path:   "/bad_request",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					if err := w.BadRequestWith(DummyError{}); err != nil {
						w.InternalServerErrorWith(err)
					}
				},
			},
			{
				Path:   "/unauthorized",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					if err := w.UnauthorizedWith(DummyError{}); err != nil {
						w.InternalServerErrorWith(err)
					}
				},
			},
			{
				Path:   "/forbidden",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					if err := w.ForbiddenWith(DummyError{}); err != nil {
						w.InternalServerErrorWith(err)
					}
				},
			},
			{
				Path:   "/not_found",
				Method: http.MethodGet,
				Handler: func(w *response.Response, r *request.Request) {
					if err := w.NotFoundWith(DummyError{}); err != nil {
						w.InternalServerErrorWith(err)
					}
				},
			},
		},
	}

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
	apiRoot.Add(ping)
	apiRoot.Add(testing.List()...)

	cfg := &config.App{
		Logger:                  logger.DefaultLogger,
		DefaultResponseRenderer: config.JSONAPIRenderer,
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
	fmt.Fprint(w, logger.DefaultLogger.SInfof(""))
	w.Flush()

	s.Run(context.Background())
}
