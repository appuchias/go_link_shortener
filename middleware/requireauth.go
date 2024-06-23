package middleware

import (
	"net/http"

	"github.com/appuchias/go_link_shortener/db"
)

// Excluded paths from the auth middleware
var excludedPaths = []string{
	"/admin/login",
}

func UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is excluded
		for _, path := range excludedPaths {
			if path == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Validate the session
		sessionid, err := db.GetKeyFromRequest(r)
		if err != nil || sessionid == "" || !db.IsSessionIDValid(sessionid) {
			http.Redirect(w, r, "/admin/login?next="+r.URL.Path, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func APIAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the session
		sessionid, err := db.GetKeyFromRequest(r)
		if err != nil || sessionid == "" || !db.IsAPIKeyValid(sessionid) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
