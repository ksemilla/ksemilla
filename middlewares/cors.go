package middlewares

import (
	"net/http"

	"ksemilla/config"

	"github.com/go-chi/cors"
)

func GetCorsHandler() func(http.Handler) http.Handler {

	config := config.Config()

	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// AllowedOrigins: []string{"https://localhost:3000", "http://localhost:3000", "http://localhost:5000"},
		AllowedOrigins: config.CORS_ALLOWED_HOSTS,
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
