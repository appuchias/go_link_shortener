package api

import (
	"net/http"
)

var APIRouter = http.NewServeMux()

func init() {
	APIRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})
}
