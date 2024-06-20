package admin

import (
	"log"
	"net/http"
)

var AdminRouter = http.NewServeMux()

func init() {
	log.Println("Admin router initialized")
	AdminRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Admin panel"))
	})
}
