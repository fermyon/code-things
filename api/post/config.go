package main

import (
	"fmt"
	"os"

	"github.com/fermyon/spin/sdk/go/config"
	"github.com/fermyon/spin/sdk/go/key_value"
)

// Config Helpers

type Config struct {
	Issuer   string
	Audience string
	JwksUrl  string
	DbUrl    string
}

func GetConfig() Config {
	domain := configGetRequired(defStore, "auth_domain")
	return Config{
		Issuer:   fmt.Sprintf("https://%v/", domain),
		Audience: configGetRequired(defStore, "auth_audience"),
		JwksUrl:  fmt.Sprintf("https://%v/.well-known/jwks.json", domain),
		DbUrl:    configGetRequired(defStore, "db_url"),
	}
}

func configGetRequired(store key_value.Store, key string) string {
	if val, err := key_value.Get(store, key); err == nil {
		return string(val)
	}
	if val, err := config.Get(key); err == nil {
		return val
	}
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	panic(fmt.Sprintf("Missing required config value: %v", key))
}
