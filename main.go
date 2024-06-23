package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/appuchias/go_link_shortener/admin"
	"github.com/appuchias/go_link_shortener/api"
	"github.com/appuchias/go_link_shortener/db"
	"github.com/appuchias/go_link_shortener/middleware"
)

const port int16 = 8000

func main() {
	// Connect to the database and close the connection when the server is stopped
	database := db.Connect()
	defer database.Close()

	// Start the server
	s := server()
	log.Fatal(s.ListenAndServe())
}

func server() *http.Server {
	router := http.NewServeMux()
	router.Handle("/admin/", middleware.UserAuth(http.StripPrefix("/admin", admin.AdminRouter)))
	router.Handle("/rest/", middleware.APIAuth(http.StripPrefix("/rest", api.APIRouter)))
	router.HandleFunc("/", http.RedirectHandler("/admin/", http.StatusMovedPermanently).ServeHTTP)
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		imgbuf, err := os.ReadFile("static/favicon.ico")
		if err != nil {
			log.Println("Favicon reading error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get last modified time
		fi, err := os.Stat("static/favicon.ico")
		if err != nil {
			log.Println("Favicon stat error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))
		w.Write(imgbuf)
	})
	router.HandleFunc("/{route}", shortenerHandler)

	middlewareStack := middleware.Chain(
		middleware.Logging,
		// middleware.IsAuthed,
		middleware.SetHeaders,
	)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middlewareStack(router),
	}
	fmt.Println("Server is running on port", port)

	return s
}

func shortenerHandler(w http.ResponseWriter, r *http.Request) {
	route := r.PathValue("route")

	if route == "" || route == "admin" {
		http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
		return
	}

	// Get the URL from the database and do the redirect
	var dst string
	dst, err := db.GetURLRedirect(route)
	if err != nil || dst == "" {
		http.Error(w, "Not found", http.StatusNotFound)
	}

	http.Redirect(w, r, dst, http.StatusMovedPermanently)
}
