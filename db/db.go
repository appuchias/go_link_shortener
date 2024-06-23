package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbPath = "shortener.db"

func init() {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables(db)

	log.Println("Database initialized successfully")
}

func createTables(db *sql.DB) {
	// Diagram made in dbdiagram.io using DBML (https://dbdiagram.io/d/Go-Link-Shortener-66775a2e5a764b3c72289b5b)

	// User
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS User (
		id_user INTEGER PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		salt TEXT NOT NULL,
		pwd TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal("Error creating User table:", err)
	}
	var user_count int
	err = db.QueryRow("SELECT COUNT(*) FROM User").Scan(&user_count)
	if err != nil {
		log.Fatal("Error counting users in User table:", err)
	}
	if user_count == 0 {
		_, err = db.Exec(`INSERT INTO User (username, salt, pwd) VALUES ("admin", "salt", "password");`)
		if err != nil {
			log.Fatal("Error inserting user into User table:", err)
		}
	}

	// Session
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Session (
		id_user INTEGER,
		valid_from INTEGER,
		valid_until INTEGER,
		key TEXT NOT NULL,
		api INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (id_user, valid_from),
		FOREIGN KEY (id_user) REFERENCES User(id_user)
	);`)
	if err != nil {
		log.Fatal("Error creating Session table:", err)
	}

	// Link
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Link (
		id_link INTEGER PRIMARY KEY,
		owner INTEGER NOT NULL,
		src TEXT NOT NULL,
		dst TEXT NOT NULL,
		is_slug BOOLEAN NOT NULL,
		FOREIGN KEY (owner) REFERENCES User(id_user)
	);`)
	if err != nil {
		log.Fatal("Error creating Link table:", err)
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS "index_link_src" ON "Link" ("src");`)
	if err != nil {
		log.Fatal("Error creating index on Link table:", err)
	}

	// Visit
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Visit (
		id_link INTEGER NOT NULL,
		datetime INTEGER NOT NULL,
		user_agent TEXT NOT NULL,
		PRIMARY KEY (id_link, datetime),
		FOREIGN KEY (id_link) REFERENCES Link(id_link)
	);`)
	if err != nil {
		log.Fatal("Error creating Visit table:", err)
	}
}

func Connect() *sql.DB {
	return db
}

func Close() {
	db.Close()
}

func GetURLRedirect(src string) (string, error) {
	var dst string
	err := db.QueryRow("SELECT dst FROM Link WHERE src = ?", src).Scan(&dst)
	if err == sql.ErrNoRows {
		// log.Println("No rows found")
		return "", nil
	} else if err != nil {
		log.Println("Error getting URL redirect from DB:", err)
		return "", err
	}
	return dst, nil
}

func IsSessionIDValid(sessionid string) bool {
	var valid_from, valid_until int64
	err := db.QueryRow("SELECT valid_from, valid_until FROM Session WHERE key = ?", sessionid).Scan(&valid_from, &valid_until)
	if err == sql.ErrNoRows {
		// log.Println("No rows found")
		return false
	} else if err != nil {
		log.Println("Error checking if sessionid is valid:", err)
		return false
	}

	if valid_from < time.Now().Unix() && valid_until > time.Now().Unix() {
		return true
	}

	return false
}
