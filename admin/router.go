package admin

import (
	"net/http"
)

var AdminRouter = http.NewServeMux()
var urlsRouter = http.NewServeMux()

func init() {
	AdminRouter.Handle("/urls/", http.StripPrefix("/urls", urlsRouter))
	AdminRouter.HandleFunc("/", adminHandler)
	AdminRouter.HandleFunc("POST /login", loginHandler)
	AdminRouter.HandleFunc("GET /login", loginHandler)
	AdminRouter.HandleFunc("GET /logout", logoutHandler)
	AdminRouter.HandleFunc("GET /password", passwordHandler)
	AdminRouter.HandleFunc("POST /password", passwordHandler)

	urlsRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})
	urlsRouter.HandleFunc("GET /edit/{id_link}", editableURLHandler)
	urlsRouter.HandleFunc("POST /add", addURLHandler)
	urlsRouter.HandleFunc("PUT /{id_link}", updateURLHandler)
	urlsRouter.HandleFunc("DELETE /{id_link}", deleteURLHandler)
}
