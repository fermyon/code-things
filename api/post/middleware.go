package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MicahParks/keyfunc"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
)

// Authorization Middleware

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

// Post Model Middleware

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
