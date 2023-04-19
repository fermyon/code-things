module github.com/fermyon/code-things/api/post

go 1.20

replace github.com/fermyon/spin/sdk/go v1.0.0 => github.com/jpflueger/spin/sdk/go v0.6.1-0.20230405131322-423d9b11be46

require (
	github.com/MicahParks/keyfunc v1.9.0
	github.com/fermyon/spin/sdk/go v1.0.0
	github.com/go-chi/chi/v5 v5.0.8
	github.com/golang-jwt/jwt/v4 v4.5.0
)
