package admin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const templatesDir = "templates/"

var AdminRouter = http.NewServeMux()

func check(err error, w http.ResponseWriter, msg string) bool {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	return err != nil
}

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
		// fmt.Fprint(w, "Login page")

		// Parse the login template
		loginTemplate, err := template.ParseFiles(templatesDir+"base.html", templatesDir+"admin/login.html")
		if check(err, w, "Error parsing login template") {
			return
		}

		// Execute the login template
		err = loginTemplate.Execute(w, struct{ Title string }{Title: "Login"})
		if check(err, w, "Error executing login template") {
			return
		}

		return
	}

	fmt.Fprint(w, "Login action")
	// // Set the sessionid cookie
	// http.SetCookie(w, &http.Cookie{
	// 	Name:  "sessionid",
	// 	Value: "",
	//	Path:   "/",
	//	HttpOnly: true, // No JS
	// //	Secure: true, // HTTPS only
	//	MaxAge: 3600,
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
