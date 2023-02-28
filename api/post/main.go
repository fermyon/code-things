package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/go-chi/chi/v5"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, req *http.Request) {
		// we need to setup the router in the spin handler
		r := chi.NewRouter()

		// r.Use(middleware.RequestID)
		// r.Use(middleware.Logger)
		// r.Use(middleware.Recoverer)

		// fmt.Println("Headers: ")
		// for k, v := range req.Header {
		// 	fmt.Printf("\t%s=%s\n", k, v)
		// }
		// fmt.Println("---")

		// mount a sub-router for posts to the appropriate component path
		// otherwise we would need to concat the prefix to all of the
		// routes which seems messier
		routePrefix := req.Header.Get("Spin-Component-Route")
		r.Mount(routePrefix, Posts())

		r.ServeHTTP(w, req)
	})
}

func main() {}

func Posts() chi.Router {
	r := chi.NewRouter()

	r.Post("/", createPost)
	r.Get("/", listPosts)
	r.Get("/{id:[0-9]+}", readPost)
	r.Put("/{id:[0-9]+}", updatePost)
	r.Delete("/{id:[0-9]+}", deletePost)

	return r
}

func notFound(w http.ResponseWriter, r *http.Request) {
}

func createPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create Post not yet implemented")
}

func readPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Read Post not yet implemented")
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	// var posts = []*Post{
	// 	{ID: "1", UserID: 100},
	// 	{ID: "2", UserID: 200},
	// 	{ID: "3", UserID: 300},
	// 	{ID: "4", UserID: 400},
	// 	{ID: "5", UserID: 500},
	// }
	fmt.Fprintln(w, "List Posts not yet implemented")
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Update Post not yet implemented")
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Delete Post not yet implemented")
}

/*
Response Rendering
*/
type Post struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"`
}
