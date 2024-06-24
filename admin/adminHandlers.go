package admin

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/appuchias/go_link_shortener/db"
)

const templatesDir = "templates/"

func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the admin template
	adminTemplate, err := template.ParseFiles(templatesDir+"base.html", templatesDir+"admin/index.html")
	if err != nil {
		log.Println("Error parsing admin template", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare the template data
	username, err := db.GetCurrentUsername(r)
	if err != nil {
		log.Println("Error getting current username", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	urls, err := db.GetUserURLs(username)
	if err != nil {
		log.Println("Error getting user URLs", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Execute the admin template
	err = adminTemplate.Execute(w, struct {
		Username string
		URLs     []db.URL
	}{Username: username, URLs: urls})
	if err != nil {
		log.Println("Error executing admin template", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Redirect authenticated users to the admin panel
		sessionid, err := db.GetKeyFromRequest(r)
		if err == nil && db.IsSessionIDValid(sessionid) {
			http.Redirect(w, r, "/admin", http.StatusFound)
			return
		}

		// Parse the login template
		loginTemplate, err := template.ParseFiles(templatesDir+"base.html", templatesDir+"admin/login.html")
		if err != nil {
			log.Println("Error parsing login template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Execute the login template
		err = loginTemplate.Execute(w, struct {
			Title string
			Next  string
		}{Title: "Login", Next: r.URL.Query().Get("next")})
		if err != nil {
			log.Println("Error executing login template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the login
	id_user, salt, hashedPwd, err := db.GetUserDetails(r.FormValue("username"))
	if err != nil {
		log.Println("Invalid credentials", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if db.HashPassword(r.FormValue("password"), salt) != hashedPwd {
		log.Println("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a new session ID lasting 7 days
	sessionid, err := db.NewSessionID(id_user, int(time.Now().Add(7*24*time.Hour).Unix()), false)
	if err != nil {
		log.Println("Error creating new session ID")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the sessionid cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionid",
		Value:    sessionid,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true, // No JS
		// Secure:   true, // HTTPS only
	})

	// Sleep randomly to prevent timing attacks
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// Redirect to the next page
	if r.URL.Query().Get("next") != "" {
		http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusFound)
	}

	// Redirect to the admin panel
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Delete the sessionid cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionid",
		Path:   "/",
		MaxAge: -1,
	})

	// Invalidate the session ID
	sessionid, err := db.GetKeyFromRequest(r)
	if err == nil {
		db.InvalidateSessionID(sessionid)
	}

	// Redirect to the login page
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func passwordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Redirect unauthenticated users to the login page
		sessionid, err := db.GetKeyFromRequest(r)
		if err != nil || !db.IsSessionIDValid(sessionid) {
			http.Redirect(w, r, "/admin/login?next=/admin/password", http.StatusFound)
			return
		}

		// Parse the password template
		passwordTemplate, err := template.ParseFiles(templatesDir+"base.html", templatesDir+"admin/password.html")
		if err != nil {
			log.Println("Error parsing password template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Execute the password template
		err = passwordTemplate.Execute(w, struct {
			Title string
			Next  string
		}{Title: "Change Password", Next: r.URL.Query().Get("next")})
		if err != nil {
			log.Println("Error executing password template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	// Redirect unauthenticated users to the login page
	sessionid, err := db.GetKeyFromRequest(r)
	if err != nil || !db.IsSessionIDValid(sessionid) {
		http.Redirect(w, r, "/admin/login?next=/admin/password", http.StatusFound)
		return
	}

	// Parse the form
	err = r.ParseForm()
	if err != nil {
		log.Println("Error parsing form")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the password
	username, err := db.GetCurrentUsername(r)
	if err != nil {
		log.Println("Error getting current username", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	id_user, salt, hashedPwd, err := db.GetUserDetails(username)
	if err != nil {
		log.Println("Error getting user details", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if db.HashPassword(r.FormValue("password"), salt) != hashedPwd {
		log.Println("Invalid credentials")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if r.FormValue("newpassword") != r.FormValue("newpassword2") {
		log.Println("Passwords don't match")
		http.Error(w, "Passwords don't match", http.StatusBadRequest)
		return
	}

	// Change the password
	err = db.ChangePassword(id_user, r.FormValue("newpassword"))
	if err != nil {
		log.Println("Error changing password", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Remove the session IDs
	db.InvalidateAllSessionIDs(id_user)

	// Redirect to the next page
	if r.URL.Query().Get("next") != "" {
		http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusFound)
	}

	// Redirect to the admin panel
	http.Redirect(w, r, "/admin", http.StatusFound)
}
