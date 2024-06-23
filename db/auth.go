package db

import (
	"crypto/sha512"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

func HashPassword(password, salt string) string {
	hash := sha512.Sum512([]byte(salt + password))

	return base64.StdEncoding.EncodeToString(hash[:])
}

// Validate the user's credentials
func ValidatePassword(password, salt, hashedPwd string) bool {
	return HashPassword(password, salt) == hashedPwd
}

// The session ID exists, is valid, and is not an API key
func IsSessionIDValid(sessionid string) bool {
	_, valid_from, valid_until, api, err := getSessionIDDetails(sessionid)
	if err != nil {
		return false
	}

	return valid_from < int(time.Now().Unix()) && valid_until > int(time.Now().Unix()) && !api
}

// The session ID exists, is valid, and is an API key
func IsAPIKeyValid(sessionid string) bool {
	_, valid_from, valid_until, api, err := getSessionIDDetails(sessionid)
	if err != nil {
		return false
	}

	return valid_from < int(time.Now().Unix()) && valid_until > int(time.Now().Unix()) && api
}

// Check if the request user is authenticated
func IsUserAuthed(r *http.Request) bool {
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		log.Println("Error getting sessionid cookie:", err)
		return false
	}

	return IsSessionIDValid(cookie.Value)
}

// Check if the API request is authenticated
func IsAPIAuthed(r *http.Request) bool {
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		log.Println("Error getting sessionid cookie:", err)
		return false
	}

	return IsAPIKeyValid(cookie.Value)
}
