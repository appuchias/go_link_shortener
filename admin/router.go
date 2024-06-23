package admin

import (
	"fmt"
	"log"
	"net/http"
)

var AdminRouter = http.NewServeMux()

func init() {
	log.Println("Admin router initialized")
	AdminRouter.HandleFunc("GET /", adminHandler)
	AdminRouter.HandleFunc("GET /login", loginHandler)
	AdminRouter.HandleFunc("POST /login", loginHandler)
	AdminRouter.HandleFunc("GET /logout", logoutHandler)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Admin panel")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprint(w, "Login page")
		return
	}

	fmt.Fprint(w, "Login action")
	// // Set the sessionid cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:  "sessionid",
	// 	Value: "",
	// })

	// // Redirect to the next page
	// if r.URL.Query().Get("next") != "" {
	// 	http.Redirect(w, r, r.URL.Query().Get("next"), http.StatusMovedPermanently)
	// }

	// // Redirect to the admin panel
	// http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logout action")

	// // Delete the sessionid cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:   "sessionid",
	// 	Path:   "/",
	// 	MaxAge: -1,
	// })

	// // Redirect to the login page
	// http.Redirect(w, r, "/admin/login", http.StatusMovedPermanently)
}
