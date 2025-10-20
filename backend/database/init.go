package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func Init() {
	var err error
	Db, err = sql.Open("sqlite3", "./backend/database/sqlite.db")
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	for _, stmt := range []string{
		// Users
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			password_hash TEXT NOT NULL
		);`,

		// Posts

		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL CHECK(length(title) > 0),
			content TEXT NOT NULL CHECK(length(content) > 0),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			likes_count INTEGER DEFAULT 0,
			dislikes_count INTEGER DEFAULT 0,
			comments_count INTEGER DEFAULT 0
		);`,

		// Categories

		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			categories TEXT NOT NULL,
			post_id INTEGER REFERENCES posts(id)
		);`,

		// Comments

		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		// Likes
		`CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL
		);`,

		// Sessions
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			username TEXT NOT NULL,
			expires_at DATETIME NOT NULL CHECK(expires_at > created_at),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	} {
		if _, err := Db.Exec(stmt); err != nil {
			log.Fatal("Table creation failed:", err)
		}
	}

	defaults := []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"}
	for _, c := range defaults {
		Db.Exec(`INSERT OR IGNORE INTO categories (categories) VALUES (?)`, c)
	}
}

// Example inserts
func InsertPost() {
	_, err := Db.Exec(`
	INSERT INTO posts (user_id, title, content, created_at, likes_count, comments_count)
	VALUES (?, ?, ?, ?, ?, ?)
	`, 1, "Cybersecurity", "Thank god.", "2025-12-5", 2, 9)
	if err != nil {
		log.Println("Skipping training insert:", err)
	}
}

func InsetComment() {
	_, err := Db.Exec(`
	INSERT INTO comments ( id ,post_id , user_id ,content ,created_at) VALUES(?,?,?,?,?)
	`, 2, 3, 5, "I love you", time.Now().Format("2001-12-3"))
	if err != nil {
		log.Fatal("error in isert:", err)
	}
}
