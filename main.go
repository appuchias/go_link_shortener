package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/appuchias/go_link_shortener/admin"
	"github.com/appuchias/go_link_shortener/middleware"
)

const port int16 = 8080

func main() {
	router := http.NewServeMux()
	router.Handle("/admin/", http.StripPrefix("/admin", admin.AdminRouter))
	router.HandleFunc("GET /{route}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("route") == "admin" {
			http.Redirect(w, r, "/admin/", http.StatusMovedPermanently)
			return
		}

		fmt.Fprint(w, r.PathValue("route"))
	})

	middlewareStack := middleware.Chain(
		middleware.Logging,
		middleware.IsAuthed,
	)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middlewareStack(router),
	}

	fmt.Println("Server is running on port", port)
	log.Fatal(s.ListenAndServe())
}
