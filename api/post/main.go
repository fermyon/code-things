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
		http.Error(res, "Forbidden: You do not have permissions to perform this action", http.StatusForbidden)
		return
	}

	if err := DbInsert(&post); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJsonResponse(res, post.ToJson())
	res.Header().Add("location", fmt.Sprintf("/api/post/%v", post.ID))
	res.WriteHeader(http.StatusCreated)
}

func listPosts(res http.ResponseWriter, req *http.Request) {
	limit, offset := getPaginationParams(req)

	if posts, err := DbReadAll(limit, offset); err == nil {
		renderJsonResponse(res, ToJson(posts))
	} else {
		http.Error(res, err.Error(), http.StatusInternalServerError)
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
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if (post == Post{}) {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}
	renderJsonResponse(res, post.ToJson())
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
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	renderJsonResponse(res, post.ToJson())
}

func deletePost(res http.ResponseWriter, req *http.Request) {
	id, err := getPostId(req)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse URL param 'id': %v", id)
		renderBadRequestResponse(res, msg)
		return
	}

	if err := DbDelete(id); err == nil {
		res.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
