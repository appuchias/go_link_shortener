package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/appuchias/go_link_shortener/admin"
	"github.com/appuchias/go_link_shortener/db"
	"github.com/appuchias/go_link_shortener/middleware"
)

const port int16 = 8000

func main() {
	// Connect to the database and close the connection when the server is stopped
	database := db.Connect()
	defer database.Close()

	// Start the server
	s := Server()
	log.Fatal(s.ListenAndServe())
}

func Server() *http.Server {
	router := http.NewServeMux()
	router.Handle("GET /admin/", http.StripPrefix("/admin", admin.AdminRouter))
	router.HandleFunc("GET /", shortenerHandler)
	router.HandleFunc("GET /{route}", shortenerHandler)

	middlewareStack := middleware.Chain(
		middleware.Logging,
		middleware.IsAuthed,
	)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middlewareStack(router),
	}
	fmt.Println("Server is running on port", port)

	return s
}

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	route := r.PathValue("route")

	if route == "" || route == "admin" {
		http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
		return
	}

	// Get the URL from the database and do the redirect
	var dst string
	dst, err = db.GetURLRedirect(route)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	if dst == "" {
		http.Error(w, "Not found", http.StatusNotFound)
	}

	http.Redirect(w, r, dst, http.StatusTemporaryRedirect)

}
