package middleware

import (
	"net/http"

	"github.com/appuchias/go_link_shortener/db"
)

// Excluded paths from the auth middleware
var excludedPaths = []string{
	"/admin/login",
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is excluded
		for _, path := range excludedPaths {
			if path == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Validate the session
		sessionid, err := db.GetSessionIDFromRequest(r)
		if err != nil || sessionid == "" || !db.IsSessionIDValid(sessionid) {
			http.Redirect(w, r, "/admin/login?next="+r.URL.Path, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
