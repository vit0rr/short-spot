package router

import "net/http"

// SetResponseTypeToJSON is a middleware sets the response type to JSON
func SetResponseTypeToJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
