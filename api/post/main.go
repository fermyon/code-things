package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/fermyon/spin/sdk/go/config"
	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/fermyon/spin/sdk/go/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/valyala/fastjson"
	"golang.org/x/exp/slices"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, req *http.Request) {
		// we need to setup the router in the spin handler
		r := chi.NewRouter()

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

func createPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Post Handler invoked")

	// parse the post
	post, err := parsePostJson(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the input
	if err = post.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// insert into the database
	err = post.dbInsert()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write the post to the response as json
	_, err = io.WriteString(w, post.toJson())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success, write status and headers
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}

func readPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting id to integer: %v", err), http.StatusInternalServerError)
		return
	}

	post, err := dbReadById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write the post to the response as json
	_, err = io.WriteString(w, post.toJson())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success, write status and headers
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
}

func listPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := dbReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// until we can find a library that serializes, we need to manually build the json array
	var sb strings.Builder
	sb.WriteString("[")
	for i := 0; i < len(posts); i++ {
		json := posts[i].toJson()
		sb.WriteString(json)
		if i != len(posts)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")

	_, err = io.WriteString(w, sb.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
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
	ID         int    `json:"id"`
	AuthorID   string `json:"author_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Data       string `json:"data"`
	Visibility string `json:"visibility"`
}

var postTypes = []string{
	"permalink-range",
	"code",
}
var postVisibilities = []string{
	"public",
}

func (post *Post) validate() error {
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

func parsePostJson(r io.ReadCloser) (Post, error) {
	var post Post
	b, err := io.ReadAll(r)
	if err != nil {
		return post, fmt.Errorf("Error reading the request: %v", err)
	}

	var p fastjson.Parser
	if val, err := p.ParseBytes(b); err != nil {
		return post, fmt.Errorf("Error parsing json: %v", err)
	} else {
		post.AuthorID = string(val.GetStringBytes("author_id"))
		post.Content = string(val.GetStringBytes("content"))
		post.Type = string(val.GetStringBytes("type"))
		post.Data = string(val.GetStringBytes("data"))
		post.Visibility = string(val.GetStringBytes("visibility"))
		return post, nil
	}
}

func (post *Post) toJson() string {
	return fmt.Sprintf(`{
		"id": %v,
		"author_id": "%v",
		"content": "%v",
		"type": "%v",
		"data": "%v",
		"visibility": "%v"}`,
		post.ID,
		post.AuthorID,
		post.Content,
		post.Type,
		post.Data,
		post.Visibility)
}

func fromRow(row []postgres.DbValue) Post {
	var post Post
	post.ID = int(row[0].GetInt32())
	post.AuthorID = row[1].GetStr()
	post.Content = row[2].GetStr()
	post.Type = row[3].GetStr()
	post.Data = row[4].GetStr()
	post.Visibility = row[5].GetStr()
	return post
}

func dbReadById(id int) (Post, error) {
	db_url, err := config.Get("db_url")
	if err != nil {
		return Post{}, fmt.Errorf("Error reading db_url from config: %s", err.Error())
	}

	statement := "SELECT id, author_id, content, type, data, visibility FROM posts WHERE id=$1"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt32(int32(id)),
	}
	rowset, err := postgres.Query(db_url, statement, params)
	if err != nil {
		return Post{}, fmt.Errorf("Error reading from database: %s", err.Error())
	}

	post := fromRow(rowset.Rows[0])
	return post, nil
}

func dbReadAll() ([]Post, error) {
	db_url, err := config.Get("db_url")
	if err != nil {
		return []Post{}, fmt.Errorf("Error reading db_url from config: %s", err.Error())
	}

	statement := "SELECT id, author_id, content, type, data, visibility FROM posts"
	rowset, err := postgres.Query(db_url, statement, []postgres.ParameterValue{})
	if err != nil {
		return []Post{}, fmt.Errorf("Error reading from database: %s", err.Error())
	}

	n := len(rowset.Rows)
	posts := make([]Post, n)
	for i := 0; i < n; i++ {
		row := rowset.Rows[i]
		posts[i] = fromRow(row)
	}

	return posts, nil
}

func (post *Post) dbInsert() error {
	db_url, err := config.Get("db_url")
	if err != nil {
		return fmt.Errorf("Error reading db_url from config: %s", err.Error())
	}

	statement := "INSERT INTO posts (author_id, content, type, data, visibility) VALUES ($1, $2, $3, $4, $5)"
	params := []postgres.ParameterValue{
		postgres.ParameterValueStr(post.AuthorID),
		postgres.ParameterValueStr(post.Content),
		postgres.ParameterValueStr(post.Type),
		postgres.ParameterValueStr(post.Content),
		postgres.ParameterValueStr(post.Visibility),
	}
	_, err = postgres.Execute(db_url, statement, params)
	if err != nil {
		return fmt.Errorf("Error inserting into database: %s", err.Error())
	}

	// this is a gross hack that will surely bite me later
	rowset, err := postgres.Query(db_url, "SELECT lastval()", []postgres.ParameterValue{})
	if err != nil {
		return fmt.Errorf("Error querying id from database: %s", err.Error())
	}
	post.ID = int(rowset.Rows[0][0].GetInt64())

	return nil
}
