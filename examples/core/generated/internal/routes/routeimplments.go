package routes

import (
	"github.com/version-1/gooo/examples/core/generated/internal/schema"
	"github.com/version-1/gooo/pkg/core/api/request"
	"github.com/version-1/gooo/pkg/core/api/response"
)

func GetUsers(res *response.Response[schema.User], req *request.Request[request.Void]) {
	// do something
}

func PostUsers(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
	// do something
}

func PatchUsersId(res *response.Response[schema.User], req *request.Request[schema.MutateUser]) {
	// do something
}

func DeleteUsersId(res *response.Response[schema.User], req *request.Request[request.Void]) {
	// do something
}

func GetUsersId(res *response.Response[schema.User], req *request.Request[request.Void]) {
	// do something
}

func GetPosts(res *response.Response[schema.Post], req *request.Request[request.Void]) {
	// do something
}

func PostPosts(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
	// do something
}

func PatchPostsId(res *response.Response[schema.Post], req *request.Request[schema.MutatePost]) {
	// do something
}

func DeletePostsId(res *response.Response[schema.Post], req *request.Request[request.Void]) {
	// do something
}

func GetPostsId(res *response.Response[schema.Post], req *request.Request[request.Void]) {
	// do something
}
