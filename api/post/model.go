package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
	"golang.org/x/exp/slices"
)

// Post model

// TODO: finish replacing fastjson with encoding/json
type Post struct {
	ID         int    `json:"id,omitempty"`         // auto-incremented by postgres
	AuthorID   string `json:"author_id,omitempty"`  // foreign key to user's id
	Content    string `json:"content,omitempty"`    // anything the poster wants to say about a piece of code they're sharing
	Type       string `json:"type,omitempty"`       // post could be a permalink, pasted code, gist, etc.
	Data       string `json:"data,omitempty"`       // actual permalink, code, gist link, etc.
	Visibility string `json:"visibility,omitempty"` // basic visibility of public, friends, etc.
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
