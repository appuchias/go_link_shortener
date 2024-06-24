package db

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = validChars[rand.Intn(len(validChars))]
	}
	return string(b)
}

func GetUserDetails(username string) (int, string, string, error) {
	var id_user int
	var salt, hashedPwd string
	err := db.QueryRow("SELECT id_user, salt, pwd FROM User WHERE username = ?", username).Scan(&id_user, &salt, &hashedPwd)
	if err != nil {
		log.Println("Error getting user ID from DB:", err)
		return 0, "", "", err
	}

	return id_user, salt, hashedPwd, nil
}

func CreateUser(username string, password string) error {
	salt := RandString(16)

	_, err := db.Exec("INSERT INTO User (username, salt, pwd) VALUES (?, ?, ?)", username, salt, HashPassword(password, salt))
	if err != nil {
		log.Println("Error creating user in DB:", err)
		return err
	}

	return nil
}

func ChangePassword(id_user int, password string) error {
	salt := RandString(16)

	_, err := db.Exec("UPDATE User SET salt = ?, pwd = ? WHERE id_user = ?", salt, HashPassword(password, salt), id_user)
	if err != nil {
		log.Println("Error changing password in DB:", err)
		return err
	}

	return nil
}

func getUsername(id_user int) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM User WHERE id_user = ?", id_user).Scan(&username)
	if err != nil {
		log.Println("Error getting username from DB:", err)
		return "", err
	}

	return username, nil
}

// Sessions

func NewSessionID(id_user int, valid_until int, api bool) (string, error) {
	// Generate a random key
	key := RandString(32)

	// Insert the session ID into the database
	_, err := db.Exec("INSERT INTO Session (id_user, valid_from, valid_until, key, api) VALUES (?, ?, ?, ?, ?)", id_user, time.Now().Unix(), valid_until, key, api)
	if err != nil {
		log.Println("Error inserting sessionid into DB:", err)
		return "", err
	}

	return key, nil
}

// Get the information of the session ID
func getSessionIDDetails(sessionid string) (int, int, int, bool, error) {
	var id_user, valid_from, valid_until int
	var api bool
	err := db.QueryRow("SELECT id_user, valid_from, valid_until, api FROM Session WHERE key = ?", sessionid).Scan(&id_user, &valid_from, &valid_until, &api)
	if err != nil {
		log.Println("Error getting sessionid details from DB:", err)
		return 0, 0, 0, false, err
	}

	return id_user, valid_from, valid_until, api, nil
}

// Get the user ID from the session ID
func GetUserIDFromSessionID(sessionid string) (int, error) {
	id_user, _, _, _, err := getSessionIDDetails(sessionid)
	if err != nil {
		return 0, err
	}

	return id_user, nil
}

// Get the username from the session ID
func GetCurrentUsername(r *http.Request) (string, error) {
	sessionid, err := GetKeyFromRequest(r)
	if err != nil {
		return "", err
	}

	id_user, err := GetUserIDFromSessionID(sessionid)
	if err != nil {
		return "", err
	}

	username, err := getUsername(id_user)
	if err != nil {
		return "", err
	}

	return username, nil
}
