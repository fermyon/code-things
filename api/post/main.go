package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/fermyon/spin/sdk/go/config"
	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/fermyon/spin/sdk/go/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/valyala/fastjson"
	"golang.org/x/exp/slices"
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

func getPaginationParams(req *http.Request) (limit int, offset int) {
	// helper function to clamp the value
	clamp := func(val int, min int, max int) int {
		if val < min {
			return min
		} else if val > max {
			return max
		} else {
			return val
		}
	}

	// get the limit from the URL
	limit_param := chi.URLParam(req, "limit")
	if limit_val, err := strconv.Atoi(limit_param); err != nil {
		// error occurred, just use a default value
		fmt.Printf("Failed to parse the limit from URL: %v", err)
		limit = 5
	} else {
		// clamp the value in case of invalid parameters (intentional or otherwise)
		limit = clamp(limit_val, 0, 25)
	}

	// get the offset from the URL
	offset_param := chi.URLParam(req, "offset")
	if offset_val, err := strconv.Atoi(offset_param); err != nil {
		// error occurred, just use a default value
		fmt.Printf("Failed to parse the offset from URL: %v", err)
		offset = 0
	} else {
		// clamp the value in case of invalid parameters (intentional or otherwise)
		// limiting this one to 10,000 because I find it unlikely that anyone will post 10k times :)
		offset = clamp(offset_val, 0, 10000)
	}

	return limit, offset
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
	db_url := getDbUrl()
	statement := "INSERT INTO posts (author_id, content, type, data, visibility) VALUES ($1, $2, $3, $4, $5)"
	params := []postgres.ParameterValue{
		postgres.ParameterValueStr(post.AuthorID),
		postgres.ParameterValueStr(post.Content),
		postgres.ParameterValueStr(post.Type),
		postgres.ParameterValueStr(post.Data),
		postgres.ParameterValueStr(post.Visibility),
	}

	_, err := postgres.Execute(db_url, statement, params)
	if err != nil {
		return fmt.Errorf("Error inserting into database: %s", err.Error())
	}

	// this is a gross hack that will surely bite me later
	rowset, err := postgres.Query(db_url, "SELECT lastval()", []postgres.ParameterValue{})
	if err != nil || len(rowset.Rows) != 1 || len(rowset.Rows[0]) != 1 {
		return fmt.Errorf("Error querying id from database: %s", err.Error())
	}

	id_val := rowset.Rows[0][0]
	if id_val.Kind() == postgres.DbValueKindInt64 {
		post.ID = int(id_val.GetInt64())
	} else {
		fmt.Printf("Failed to populate created post's identifier, invalid kind returned from database: %v", id_val.Kind())
	}

	return nil
}

func DbReadById(id int) (Post, error) {
	db_url := getDbUrl()
	statement := "SELECT id, author_id, content, type, data, visibility FROM posts WHERE id=$1"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt32(int32(id)),
	}

	rowset, err := postgres.Query(db_url, statement, params)
	if err != nil {
		return Post{}, fmt.Errorf("Error reading from database: %s", err.Error())
	}

	if rowset.Rows == nil || len(rowset.Rows) == 0 {
		return Post{}, nil
	} else {
		return fromRow(rowset.Rows[0])
	}
}

func DbReadAll(limit int, offset int) ([]Post, error) {
	db_url := getDbUrl()
	statement := "SELECT id, author_id, content, type, data, visibility FROM posts ORDER BY id LIMIT $1 OFFSET $2"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt64(int64(limit)),
		postgres.ParameterValueInt64(int64(offset)),
	}
	rowset, err := postgres.Query(db_url, statement, params)
	if err != nil {
		return []Post{}, fmt.Errorf("Error reading from database: %s", err.Error())
	}

	posts := make([]Post, len(rowset.Rows))
	for i, row := range rowset.Rows {
		if post, err := fromRow(row); err != nil {
			return []Post{}, err
		} else {
			posts[i] = post
		}
	}

	return posts, nil
}

func DbUpdate(post Post) error {
	db_url := getDbUrl()
	statement := "UPDATE posts SET author_id=$1, content=$2, type=$3, data=$4, visibility=$5 WHERE id=$6"
	params := []postgres.ParameterValue{
		postgres.ParameterValueStr(post.AuthorID),
		postgres.ParameterValueStr(post.Content),
		postgres.ParameterValueStr(post.Type),
		postgres.ParameterValueStr(post.Data),
		postgres.ParameterValueStr(post.Visibility),
		postgres.ParameterValueInt32(int32(post.ID)),
	}

	_, err := postgres.Execute(db_url, statement, params)
	if err != nil {
		return fmt.Errorf("Error updating database: %s", err.Error())
	}

	return nil
}

func DbDelete(id int) error {
	db_url := getDbUrl()
	statement := "DELETE FROM posts WHERE id=$1"
	params := []postgres.ParameterValue{
		postgres.ParameterValueInt32(int32(id)),
	}

	_, err := postgres.Execute(db_url, statement, params)
	return err
}

// Database Helper Functions

func assertValueKind(row []postgres.DbValue, col int, expected postgres.DbValueKind) (postgres.DbValue, error) {
	if row[col].Kind() != expected {
		return postgres.DbValue{}, fmt.Errorf("Expected column %v to be %v kind but got %v\n", col, expected, row[col].Kind())
	}
	return row[col], nil
}

func fromRow(row []postgres.DbValue) (Post, error) {
	var post Post

	if val, err := assertValueKind(row, 0, postgres.DbValueKindInt32); err == nil {
		post.ID = int(val.GetInt32())
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 1, postgres.DbValueKindStr); err == nil {
		post.AuthorID = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 2, postgres.DbValueKindStr); err == nil {
		post.Content = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 3, postgres.DbValueKindStr); err == nil {
		post.Type = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 4, postgres.DbValueKindStr); err == nil {
		post.Data = val.GetStr()
	} else {
		return post, err
	}

	if val, err := assertValueKind(row, 5, postgres.DbValueKindStr); err == nil {
		post.Visibility = val.GetStr()
	} else {
		return post, err
	}

	return post, nil
}

// Config Helpers

func configGetRequired(key string) string {
	if val, err := config.Get(key); err != nil {
		panic(fmt.Sprintf("Missing required config item 'jwks_uri': %v", err))
	} else {
		return val
	}
}

func getIssuer() string {
	domain := configGetRequired("auth_domain")
	return fmt.Sprintf("https://%v/", domain)
}

func getAudience() string {
	return configGetRequired("auth_audience")
}

func getJwksUri() string {
	domain := configGetRequired("auth_domain")
	return fmt.Sprintf("https://%v/.well-known/jwks.json", domain)
}

func getDbUrl() string {
	return configGetRequired("db_url")
}

// Authorization Helpers

type claimsCtxKey struct{}

func TokenVerification(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// ensure RS256 was used to sign the token
		parserOpts := []request.ParseFromRequestOption{
			request.WithParser(jwt.NewParser(jwt.WithValidMethods([]string{
				jwt.SigningMethodRS256.Alg(),
			}))),
		}
		token, err := request.ParseFromRequest(req, request.OAuth2Extractor, fetchAuthSigningKey, parserOpts...)
		if err != nil {
			if errors.Is(err, jwt.ValidationError{}) {
				// token parsed but was invalid
				http.Error(res, err.Error(), http.StatusUnauthorized)
			} else {
				// unable to parse or verify signing
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if !token.Valid {
			http.Error(res, "token not valid", http.StatusUnauthorized)
			return
		}

		fmt.Printf("Claims: %v\n", claims)

		if !claims.VerifyIssuer(getIssuer(), true) {
			fmt.Printf("Expected issuer %v but got %v", getIssuer(), claims["iss"])
			http.Error(res, jwt.ErrTokenInvalidIssuer.Error(), http.StatusUnauthorized)
			return
		}

		if !claims.VerifyAudience(getAudience(), true) {
			http.Error(res, jwt.ErrTokenInvalidAudience.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), claimsCtxKey{}, claims)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func fetchAuthSigningKey(t *jwt.Token) (interface{}, error) {
	jwksUri := getJwksUri()
	println("Fetching auth signing key from: %v", jwksUri)
	if jwks, err := keyfunc.Get(jwksUri, keyfunc.Options{
		Client: NewHttpClient(),
	}); err != nil {
		println("Failed to fetch auth signing key: %v", err)
		return nil, err
	} else {
		println("Successfully retrieved and parsed JWKS")
		return jwks.Keyfunc(t)
	}
}

// HTTP Helpers

type spinRoundTripper struct{}

func (t spinRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return spinhttp.Send(req)
}

func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: spinRoundTripper{},
		Timeout:   time.Duration(5) * time.Second,
	}
}
