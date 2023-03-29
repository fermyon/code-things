package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

func main() {}

func init() {
	spinhttp.Handle(func(res http.ResponseWriter, req *http.Request) {
		// we need to setup the router inside spin handler
		router := chi.NewRouter()

		router.Use(TokenVerification)

		// mount our routes using the prefix
		routePrefix := req.Header.Get("Spin-Component-Route")
		router.Mount(routePrefix, PostRouter())

		// hand the request/response off to chi
		router.ServeHTTP(res, req)
	})
}

// TODO: create wrapper function to handle errors?
func PostRouter() chi.Router {
	posts := chi.NewRouter()
	idParamPattern := fmt.Sprintf("/{%v:[0-9]+}", postIdCtxKey)
	posts.Use(PostCtx)
	posts.Post("/", createPost)
	posts.Get("/", listPosts)
	posts.Get(idParamPattern, readPost)
	posts.Put(idParamPattern, updatePost)
	posts.Delete(idParamPattern, deletePost)
	return posts
}

func createPost(res http.ResponseWriter, req *http.Request) {
	post := req.Context().Value(postCtxKey{}).(Post)
	claims := req.Context().Value(claimsCtxKey{}).(jwt.MapClaims)

	if claims["sub"] != post.AuthorID {
		renderForbiddenResponse(res)
		return
	}

	if err := DbInsert(&post); err != nil {
		renderErrorResponse(res, err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	res.Header().Add("location", fmt.Sprintf("/api/post/%v", post.ID))
	renderJsonResponse(res, post)
}

func listPosts(res http.ResponseWriter, req *http.Request) {
	limit, offset := getPaginationParams(req)

	if posts, err := DbReadAll(limit, offset); err != nil {
		renderErrorResponse(res, err)
	} else {
		renderJsonResponse(res, posts)
	}
}

func readPost(res http.ResponseWriter, req *http.Request) {
	id, err := getPostId(req)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse URL param 'id': %v", id)
		renderBadRequestResponse(res, msg)
		return
	}

	post, err := DbReadById(id)
	if err != nil {
		renderErrorResponse(res, err)
		return
	}
	if (post == Post{}) {
		renderNotFound(res)
		return
	}

	renderJsonResponse(res, post)
}

func updatePost(res http.ResponseWriter, req *http.Request) {
	post := req.Context().Value(postCtxKey{}).(Post)

	if id, err := getPostId(req); err != nil {
		msg := fmt.Sprintf("Failed to parse URL param 'id': %v", id)
		renderBadRequestResponse(res, msg)
		return
	} else {
		post.ID = id
	}

	if err := DbUpdate(post); err != nil {
		renderErrorResponse(res, err)
	}
	renderJsonResponse(res, post)
}

func deletePost(res http.ResponseWriter, req *http.Request) {
	id, err := getPostId(req)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse URL param 'id': %v", id)
		renderBadRequestResponse(res, msg)
		return
	}

	if err := DbDelete(id); err != nil {
		renderErrorResponse(res, err)
	}
	res.WriteHeader(http.StatusNoContent)
}
