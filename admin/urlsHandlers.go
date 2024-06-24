package admin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/appuchias/go_link_shortener/db"
)

const shortURLLength = 6

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

	editableTemplate, err := template.ParseFiles(templatesDir + "admin/editable.html")
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
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get the fields from the form
	src := r.FormValue("shorturl")
	dst := r.FormValue("longurl")
	isCustom := true

	if src == "" {
		isCustom = false
		src = db.RandString(shortURLLength)
	}

	sessionid, err := db.GetKeyFromRequest(r)
	if err != nil {
		log.Println("Error getting session ID", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	id_user, err := db.GetUserIDFromSessionID(sessionid)
	if err != nil {
		log.Println("Error getting user ID from session ID", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add the URL to the database
	err = db.AddURL(id_user, src, dst, isCustom)
	if err != nil {
		log.Println("Error adding URL to the database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/", http.StatusFound)
}

func updateURLHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if the user is the owner of the URL
	url, err := db.GetURLDetails(id_link)
	if err != nil {
		log.Println("Error getting URL details", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	sessionid, err := db.GetKeyFromRequest(r)
	if err != nil {
		log.Println("Error getting session ID", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	id_user, err := db.GetUserIDFromSessionID(sessionid)
	if err != nil {
		log.Println("Error getting user ID from session ID", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if id_user != url.Owner {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Parse the form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get the fields from the form
	src := r.FormValue("shorturl")
	dst := r.FormValue("longurl")
	isCustom := src != url.Src || url.IsCustom

	if src == "" {
		isCustom = false
		src = db.RandString(shortURLLength)
	}

	// Update the URL in the database
	err = db.UpdateURL(id_link, src, dst, isCustom)
	if err != nil {
		log.Println("Error updating URL in the database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<script>location.reload();</script>")
}

func deleteURLHandler(w http.ResponseWriter, r *http.Request) {
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

	// Check if the user is the owner of the URL
	url, err := db.GetURLDetails(id_link)
	if err != nil {
		log.Println("Error getting URL details", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	sessionid, err := db.GetKeyFromRequest(r)
	if err != nil {
		log.Println("Error getting session ID", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	id_user, err := db.GetUserIDFromSessionID(sessionid)
	if err != nil {
		log.Println("Error getting user ID from session ID", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if id_user != url.Owner {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = db.DeleteURL(id_link)
	if err != nil {
		log.Println("Error deleting URL", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<script>location.reload();</script>")
}
