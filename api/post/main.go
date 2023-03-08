package main

import (
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/go-chi/chi/v5"
)

func main() {}

func init() {
	spinhttp.Handle(func(res http.ResponseWriter, req *http.Request) {
		// we need to setup the router inside spin handler
		router := chi.NewRouter()

		// mount our routes using the prefix
		routePrefix := req.Header.Get("Spin-Component-Route")
		router.Mount(routePrefix, Post())

		// hand the request/response off to chi
		router.ServeHTTP(res, req)
	})
}

func PostRouter() chi.Router {
	posts := chi.NewRouter()
	posts.Post("/", createPost)
	posts.Get("/", listPosts)
	posts.Get("/{id:[0-9]+}", readPost)
	posts.Put("/{id:[0-9]+}", updatePost)
	posts.Delete("/{id:[0-9]+}", deletePost)
	return posts
}

func createPost(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Create post not yet implemented", http.StatusNotImplemented)
}

func listPosts(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "List posts not yet implemented", http.StatusNotImplemented)
}

func readPost(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Read post not yet implemented", http.StatusNotImplemented)
}

func updatePost(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Update post not yet implemented", http.StatusNotImplemented)
}

func deletePost(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Delete post not yet implemented", http.StatusNotImplemented)
}

type Post struct {
	ID         int    // auto-incremented by postgres
	AuthorID   string // foreign key to user's id
	Content    string // anything the poster wants to say about a piece of code they're sharing
	Type       string // post could be a permalink, pasted code, gist, etc.
	Data       string // actual permalink, code, gist link, etc.
	Visibility string // basic visibility of public, friends, etc.
}

// enumerated values for type
var postTypes = []string{
	"permalink-range",
	"code",
}

// enumerated values for visibility
var postVisibilities = []string{
	"public",
}
