package middleware

import (
	"net/http"
)

func SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add headers to the response
		w.Header().Set("Server", "AppuServer")

		next.ServeHTTP(w, r)
	})
}
