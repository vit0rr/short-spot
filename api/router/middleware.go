package router

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/vit0rr/short-spot/pkg/log"
)

// SetResponseTypeToJSON is a middleware sets the response type to JSON
func SetResponseTypeToJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware is a middleware that checks if the request has a valid auth token
func AuthMiddleware(next http.Handler) http.Handler {
	godotenv.Load()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token != os.Getenv("AUTH_TOKEN") {
			log.Error(r.Context(), "Unauthorized. Please provide a valid auth token")
			http.Error(w, "Unauthorized. Please provide a valid auth token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CorsMiddleware is a middleware that sets the CORS headers
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
