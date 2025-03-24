package routes

// This is a generated file. DO NOT EDIT manually.
import (
	"github.com/version-1/gooo/examples/core/generated/internal/schema"
	"github.com/version-1/gooo/pkg/core/api/request"
	"github.com/version-1/gooo/pkg/core/api/response"
	"github.com/version-1/gooo/pkg/core/api/route"
)

func Routes() []route.HandlerInterface {
	routes := route.GroupHandler{
		Path: "/users",
		Handlers: []route.HandlerInterface{
			route.JSON[request.Void, schema.User]().Get("/users", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[schema.MutateUser, schema.User]().Post("/users", func(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
				// do something
			}),
			route.JSON[schema.MutateUser, schema.User]().Patch("/users/{id}", func(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
				// do something
			}),
			route.JSON[request.Void, schema.User]().Delete("/users/{id}", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[request.Void, schema.User]().Get("/users/{id}", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[request.Void, schema.Post]().Get("/posts", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[schema.MutatePost, schema.Post]().Post("/posts", func(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
				// do something
			}),
			route.JSON[schema.MutatePost, schema.Post]().Patch("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
				// do something
			}),
			route.JSON[request.Void, schema.Post]().Delete("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
				// do something
			}),
			route.JSON[request.Void, schema.Post]().Get("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
				// do something
			}),
		},
	}

	return routes.Children()
}
