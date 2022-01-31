package config

import (
	"os"
	"strings"
)

type RawConfig struct {
	APP_SECRET_KEY     string
	CORS_ALLOWED_HOSTS []string
	MONGODB_URI        string
}

func Config() *RawConfig {

	allowed_hosts := strings.Split(os.Getenv("CORS_ALLOWED_HOSTS"), ",")

	return &RawConfig{
		APP_SECRET_KEY:     os.Getenv("APP_SECRET_KEY"),
		CORS_ALLOWED_HOSTS: allowed_hosts,
		MONGODB_URI:        os.Getenv("MONGODB_URI"),
	}
}
