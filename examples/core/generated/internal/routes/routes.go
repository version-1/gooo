package routes

// This is a generated file. DO NOT EDIT manually.
import (
	"github.com/version-1/gooo/examples/core/generated/internal/schema"
	"github.com/version-1/gooo/pkg/core/api/request"
	"github.com/version-1/gooo/pkg/core/api/response"
	"github.com/version-1/gooo/pkg/core/api/route"
)

func PostUsersHandler() route.HandlerInterface {
	return route.JSON[schema.MutateUser, schema.User]().Get("/users", func(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
		PostUsers(res, req)
	})
}

func GetUsersHandler() route.HandlerInterface {
	return route.JSON[request.Void, schema.User]().Get("/users", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
		GetUsers(res, req)
	})
}

func GetUsersIdHandler() route.HandlerInterface {
	return route.JSON[request.Void, schema.User]().Get("/users/{id}", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
		GetUsersId(res, req)
	})
}

func PatchUsersIdHandler() route.HandlerInterface {
	return route.JSON[schema.MutateUser, schema.User]().Get("/users/{id}", func(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
		PatchUsersId(res, req)
	})
}

func DeleteUsersIdHandler() route.HandlerInterface {
	return route.JSON[request.Void, schema.User]().Get("/users/{id}", func(res *response.Response[schema.User], req *request.Request[request.Void]) {
		DeleteUsersId(res, req)
	})
}

func GetPostsHandler() route.HandlerInterface {
	return route.JSON[request.Void, schema.Post]().Get("/posts", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
		GetPosts(res, req)
	})
}

func PostPostsHandler() route.HandlerInterface {
	return route.JSON[schema.MutatePost, schema.Post]().Get("/posts", func(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
		PostPosts(res, req)
	})
}

func GetPostsIdHandler() route.HandlerInterface {
	return route.JSON[request.Void, schema.Post]().Get("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
		GetPostsId(res, req)
	})
}

func PatchPostsIdHandler() route.HandlerInterface {
	return route.JSON[schema.MutatePost, schema.Post]().Get("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
		PatchPostsId(res, req)
	})
}

func DeletePostsIdHandler() route.HandlerInterface {
	return route.JSON[request.Void, schema.Post]().Get("/posts/{id}", func(res *response.Response[schema.Post], req *request.Request[request.Void]) {
		DeletePostsId(res, req)
	})
}

func Routes() []route.HandlerInterface {
	routes := route.GroupHandler{
		Path: "/users",

		Handlers: []route.HandlerInterface{
			PostUsersHandler(),
			GetUsersHandler(),
			GetUsersIdHandler(),
			PatchUsersIdHandler(),
			DeleteUsersIdHandler(),
			GetPostsHandler(),
			PostPostsHandler(),
			GetPostsIdHandler(),
			PatchPostsIdHandler(),
			DeletePostsIdHandler(),
		},
	}

	return routes.Children()
}
