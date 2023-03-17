package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/go-chi/chi/v5"
	"github.com/valyala/fastjson"
	"golang.org/x/exp/slices"
)

func main() {}

func init() {
	spinhttp.Handle(func(res http.ResponseWriter, req *http.Request) {
		// we need to setup the router inside spin handler
		router := chi.NewRouter()

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

var postIdCtxKey string = "id"

type postCtxKey struct{}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var post Post
		var err error

		if req.ContentLength > 0 && req.Header.Get("Content-Type") == "application/json" {
			if post, err = ParseJsonPost(req.Body); err != nil {
				// parsing failed end the request here
				msg := fmt.Sprintf("Failed to parse the post from request body: %v\n", err)
				renderBadRequestResponse(res, msg)
				return
			}
			if err = post.Validate(); err != nil {
				msg := fmt.Sprintf("Request body failed validation: %v\n", err)
				renderBadRequestResponse(res, msg)
				return
			}
		}

		ctx := context.WithValue(req.Context(), postCtxKey{}, post)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func getPostId(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, postIdCtxKey)
	return strconv.Atoi(idStr)
}

func createPost(res http.ResponseWriter, req *http.Request) {
	post := req.Context().Value(postCtxKey{}).(Post)

	err := DbInsert(&post)
	if err == nil {
		renderJsonResponse(res, post.ToJson())
		res.Header().Add("location", fmt.Sprintf("/api/post/%v", post.ID))
		res.WriteHeader(http.StatusCreated)
	} else {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func listPosts(res http.ResponseWriter, req *http.Request) {
	if posts, err := DbReadAll(); err == nil {
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

func renderJsonResponse(res http.ResponseWriter, json string) {
	// write the post to the response as json
	if _, err := io.WriteString(res, json); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	} else {
		res.WriteHeader(http.StatusOK)
		res.Header().Add("Content-Type", "application/json")
	}
}

func renderBadRequestResponse(res http.ResponseWriter, msg string) {
	fmt.Print(msg)
	http.Error(res, msg, http.StatusBadRequest)
}

// Post model

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

func ParseJsonPost(r io.ReadCloser) (Post, error) {
	var post Post

	// read the request bytes into []byte
	b, err := io.ReadAll(r)
	if err != nil {
		return post, fmt.Errorf("Error reading the request: %v", err)
	}

	// parse the []byte array
	var p fastjson.Parser
	if val, err := p.ParseBytes(b); err != nil {
		return post, fmt.Errorf("Error parsing json: %v", err)
	} else {
		post.ID = val.GetInt("id")
		post.AuthorID = string(val.GetStringBytes("author_id"))
		post.Content = string(val.GetStringBytes("content"))
		post.Type = string(val.GetStringBytes("type"))
		post.Data = string(val.GetStringBytes("data"))
		post.Visibility = string(val.GetStringBytes("visibility"))
		return post, nil
	}
}

func (post *Post) Validate() error {
	var errs []string
	if post.AuthorID == "" {
		errs = append(errs, "'author_id' is required")
	}
	if post.Content == "" {
		errs = append(errs, "'content' is required")
	}
	if !slices.Contains(postTypes, post.Type) {
		errs = append(errs, fmt.Sprintf("'type' must be one of [%v]", postTypes))
	}
	if post.Data == "" {
		errs = append(errs, "'data' is required")
	}
	if !slices.Contains(postVisibilities, post.Visibility) {
		errs = append(errs, fmt.Sprintf("'visibility' must be one of [%v]", postTypes))
	}

	if len(errs) > 0 {
		return fmt.Errorf("Request failed validation:\n%s", strings.Join(errs, "\n"))
	} else {
		return nil
	}
}

func (post *Post) ToJson() string {
	return fmt.Sprintf(`{
"id": %v,
"author_id": "%v",
"content": "%v",
"type": "%v",
"data": "%v",
"visibility": "%v"
}`,
		post.ID,
		post.AuthorID,
		post.Content,
		post.Type,
		post.Data,
		post.Visibility)
}

func ToJson(posts []Post) string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i := 0; i < len(posts); i++ {
		sb.WriteString(posts[i].ToJson())
		if i != len(posts)-1 {
			sb.WriteRune(',')
		}
	}
	sb.WriteRune(']')
	return sb.String()
}

// Database Operations

func DbInsert(post *Post) error {
	//TODO: implement
	post.ID = 1
	return nil
}

func DbReadAll() ([]Post, error) {
	//TODO: implement
	return []Post{}, nil
}

func DbReadById(id int) (Post, error) {
	//TODO: implement
	return Post{
		ID: id,
	}, nil
}

func DbUpdate(post Post) error {
	//TODO: implement
	return nil
}

func DbDelete(id int) error {
	//TODO: implement
	return nil
}
