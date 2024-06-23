package admin

import (
	"log"
	"net/http"
)

var AdminRouter = http.NewServeMux()

func init() {
	log.Println("Admin router initialized")
	AdminRouter.HandleFunc("GET /", adminHandler)
	AdminRouter.HandleFunc("POST /login", loginHandler)
	AdminRouter.HandleFunc("GET /login", loginHandler)
	AdminRouter.HandleFunc("GET /logout", logoutHandler)
}
