package middleware

import (
	"net/http"
)

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated
		// If not, redirect to the login page
		// log.Println("Checking if user is authenticated")
		next.ServeHTTP(w, r)
	})
}
