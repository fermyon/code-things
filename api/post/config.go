package main

import (
	"fmt"

	"github.com/fermyon/spin/sdk/go/config"
)

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
