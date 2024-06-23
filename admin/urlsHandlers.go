package admin

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/appuchias/go_link_shortener/db"
)

func checkHtmxHeader(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func editableURLHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHtmxHeader(r) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}

	// Get the URL ID from the query string as int
	id_link, err := strconv.Atoi(r.PathValue("id_link"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	editableTemplate, err := template.ParseFiles(templatesDir + "urls/editable.html")
	if err != nil {
		log.Println("Error parsing editable template", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	url, err := db.GetURLDetails(id_link)
	if err != nil {
		log.Println("Error getting URL details", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = editableTemplate.Execute(w, url)
	if err != nil {
		log.Println("Error executing editable template", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func addURLHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func updateURLHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func deleteURLHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
