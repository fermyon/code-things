package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/fermyon/spin/sdk/go/key_value"
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
		token, err := request.ParseFromRequest(req, request.OAuth2Extractor, cachedJwksKeyfunc, parserOpts...)
		if err != nil {
			if errors.Is(err, jwt.ValidationError{}) {
				// token parsed but was invalid
				renderUnauthorized(res, err)
			} else {
				// unable to parse or verify signing
				renderErrorResponse(res, err)
			}
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if !token.Valid {
			renderUnauthorized(res, fmt.Errorf("token not valid"))
			return
		}

		if !claims.VerifyIssuer(cfg.Issuer, true) {
			renderUnauthorized(res, jwt.ErrTokenInvalidIssuer)
			return
		}

		if !claims.VerifyAudience(cfg.Audience, true) {
			renderUnauthorized(res, jwt.ErrTokenInvalidAudience)
			return
		}

		ctx := context.WithValue(req.Context(), claimsCtxKey{}, claims)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func cachedJwksKeyfunc(t *jwt.Token) (interface{}, error) {
	if jwks, err := getCachedJwks(); err != nil {
		fmt.Println("Failed to get jwks from cache: ", err)
	} else {
		return jwks.Keyfunc(t)
	}
	if jwks, err := keyfunc.Get(cfg.JwksUrl, keyfunc.Options{
		Client: NewHttpClient(),
	}); err != nil {
		fmt.Println("Failed to get jwks from url: ", err)
		return nil, err
	} else {
		return jwks.Keyfunc(t)
	}
}

func getCachedJwks() (*keyfunc.JWKS, error) {
	if data, err := key_value.Get(defStore, "jwks_ttl"); err != nil {
		return nil, fmt.Errorf("Failed to get jwks_ttl from store: %v", err)
	} else {
		jwks_ttl := int64(binary.LittleEndian.Uint64(data))
		if jwks_ttl > time.Now().UTC().Unix() {
			if data, err := key_value.Get(defStore, "jwks"); err != nil {
				return nil, fmt.Errorf("Failed to get jwks from store: %v", err)
			} else {
				if jwks, err := keyfunc.NewJSON(data); err != nil {
					return nil, fmt.Errorf("Failed to parse jwks: %v", err)
				} else {
					return jwks, nil
				}
			}
		} else {
			return nil, fmt.Errorf("jwks is expired")
		}
	}
}

func setCachedJwks(jwks *keyfunc.JWKS) {
	data, err := json.Marshal(jwks)
	if err != nil {
		fmt.Println("Failed to marshal jwks: ", err)
		return
	}
	err = key_value.Set(defStore, "jwks", data)
	if err != nil {
		fmt.Println("Failed to set jwks in store: ", err)
		return
	}
	// set the ttl for the jwks key
	jwks_ttl := uint64(time.Now().UTC().Add(24 * time.Hour).Unix())
	jwks_data := make([]byte, 8)
	binary.LittleEndian.PutUint64(jwks_data, jwks_ttl)
	if err := key_value.Set(defStore, "jwks_ttl", jwks_data); err != nil {
		fmt.Println("Failed to set jwks_ttl in store: ", err)
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
			if post, err = DecodePost(req.Body); err != nil {
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
