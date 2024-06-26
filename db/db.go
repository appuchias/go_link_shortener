package db

import (
	"database/sql"
	"log"

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
	var err error
	// Diagram made in dbdiagram.io using DBML (https://dbdiagram.io/d/Go-Link-Shortener-66775a2e5a764b3c72289b5b)

	// User
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS User (
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
		err = CreateUser("admin", "password")
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
		api BOOLEAN NOT NULL DEFAULT 0,
		PRIMARY KEY (id_user, valid_from),
		FOREIGN KEY (id_user) REFERENCES User(id_user) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal("Error creating Session table:", err)
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS "index_session_key" ON "Session" ("key");`)
	if err != nil {
		log.Fatal("Error creating index on Session table:", err)
	}

	// Link
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Link (
		id_link INTEGER PRIMARY KEY,
		owner INTEGER NOT NULL,
		src TEXT NOT NULL,
		dst TEXT NOT NULL,
		is_custom BOOLEAN NOT NULL,
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
		user_agent TEXT,
		PRIMARY KEY (id_link, datetime),
		FOREIGN KEY (id_link) REFERENCES Link(id_link) ON DELETE CASCADE
	);`)
	if err != nil {
		log.Fatal("Error creating Visit table:", err)
	}

	// SQLite options
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Error enabling foreign keys in SQLite:", err)
	}
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Fatal("Error setting journal mode in SQLite:", err)
	}
	_, err = db.Exec("PRAGMA synchronous = NORMAL;")
	if err != nil {
		log.Fatal("Error setting synchronous mode in SQLite:", err)
	}
}

func Connect() *sql.DB {
	return db
}

func Close() {
	db.Close()
}
