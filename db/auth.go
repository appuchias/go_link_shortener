package db

import (
	"crypto/sha512"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

func GetKeyFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("sessionid")
	if err == nil {
		return cookie.Value, nil
	}

	header := r.Header.Get("X-Api-Key")
	if header != "" {
		return header, nil
	}

	return "", err
}

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

func InvalidateSessionID(sessionid string) error {
	_, err := db.Exec("UPDATE Session SET valid_until = ? WHERE key = ?", int(time.Now().Unix()), sessionid)
	if err != nil {
		log.Println("Error invalidating sessionid from DB:", err)
		return err
	}

	return nil
}
