package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Post model

type Post struct {
	ID         int            `json:"id,omitempty"`         // auto-incremented by postgres
	AuthorID   string         `json:"author_id,omitempty"`  // foreign key to user's id
	Content    string         `json:"content,omitempty"`    // anything the poster wants to say about a piece of code they're sharing
	Type       PostType       `json:"type,omitempty"`       // post could be a permalink, pasted code, gist, etc.
	Data       string         `json:"data,omitempty"`       // actual permalink, code, gist link, etc.
	Visibility PostVisibility `json:"visibility,omitempty"` // basic visibility of public, followers, etc.
}

func (p Post) Validate() error {
	var errs []error

	if p.AuthorID == "" {
		errs = append(errs, fmt.Errorf("field 'author_id' is required"))
	}

	if p.Content == "" {
		errs = append(errs, fmt.Errorf("field 'content' is required"))
	}

	if p.Type == 0 {
		errs = append(errs, fmt.Errorf("field 'type' contains unknown value"))
	}

	if p.Data == "" {
		errs = append(errs, fmt.Errorf("field 'data' is required"))
	}

	if p.Visibility == 0 {
		errs = append(errs, fmt.Errorf("field 'visibility' contains unknown value"))
	}

	return errors.Join(errs...)
}

// Post Type enum

type PostType uint8

const (
	PostTypePermalinkRange PostType = iota + 1
)

var (
	PostType_name = map[PostType]string{
		PostTypePermalinkRange: "permalink-range",
	}
	PostType_value = map[string]PostType{
		"permalink-range": PostTypePermalinkRange,
	}
)

func (t PostType) String() string {
	return PostType_name[t]
}

func ParsePostType(s string) (PostType, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := PostType_value[s]
	if !ok {
		return PostType(0), fmt.Errorf("%q is not a valid post type", s)
	}
	return PostType(value), nil
}

func (t PostType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *PostType) UnmarshalJSON(data []byte) (err error) {
	var postType string
	if err := json.Unmarshal(data, &postType); err != nil {
		return err
	}
	if *t, err = ParsePostType(postType); err != nil {
		return err
	}
	return nil
}

// Post Visibility enum

type PostVisibility uint8

const (
	PostVisibilityPublic PostVisibility = iota + 1
	PostVisibilityFollowers
)

var (
	PostVisibility_name = map[PostVisibility]string{
		PostVisibilityPublic:    "public",
		PostVisibilityFollowers: "followers",
	}
	PostVisibility_value = map[string]PostVisibility{
		"public":    PostVisibilityPublic,
		"followers": PostVisibilityFollowers,
	}
)

func (v PostVisibility) String() string {
	return PostVisibility_name[v]
}

func ParsePostVisibility(s string) (PostVisibility, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := PostVisibility_value[s]
	if !ok {
		return PostVisibility(0), fmt.Errorf("%q is not a valid post visibility", s)
	}
	return PostVisibility(value), nil
}

func (v PostVisibility) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *PostVisibility) UnmarshalJSON(data []byte) (err error) {
	var visibility string
	if err := json.Unmarshal(data, &visibility); err != nil {
		return err
	}
	if *v, err = ParsePostVisibility(visibility); err != nil {
		return err
	}
	return nil
}

// JSON helpers

func DecodePost(r io.ReadCloser) (Post, error) {
	decoder := json.NewDecoder(r)
	var p Post
	err := decoder.Decode(&p)
	return p, err
}
