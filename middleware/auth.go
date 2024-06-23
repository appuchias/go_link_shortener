package middleware

import (
	"net/http"

	"github.com/appuchias/go_link_shortener/db"
)

// Excluded paths from the auth middleware
var excludedPaths = []string{
	"/admin/login",
}

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is excluded
		for _, path := range excludedPaths {
			if path == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Check the sessionid cookie
		sessionid, err := r.Cookie("sessionid")
		if err != nil || !db.IsSessionIDValid(sessionid.Value) {
			http.Redirect(w, r, "/admin/login?next="+r.URL.Path, http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}
