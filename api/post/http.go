package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	spinhttp "github.com/fermyon/spin/sdk/go/http"
	"github.com/go-chi/chi/v5"
)

// HTTP Response Helpers

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

// Pagination Helpers

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

// HTTP Client
//TODO: this could be contributed to spin go sdk

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
