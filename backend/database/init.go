package database

import (
	"database/sql"
	"log"

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
			username TEXT NOT NULL,
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
			name TEXT UNIQUE NOT NULL
		);`,

		// Comments

		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			likes_count INTEGER DEFAULT 0,
			dislikes_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER NOT NULL REFERENCES posts(id),
			category_id INTEGER NOT NULL REFERENCES categories(id),
			UNIQUE(post_id, category_id)
		);`,

		// Comment likes
		`CREATE TABLE IF NOT EXISTS comment_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			comment_id INTEGER NOT NULL REFERENCES comments(id),
			UNIQUE(username, comment_id)
		);`,

		// Comment dislikes
		`CREATE TABLE IF NOT EXISTS comment_dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			comment_id INTEGER NOT NULL REFERENCES comments(id),
			UNIQUE(username, comment_id)
		);`,

		// Likes
		`CREATE TABLE IF NOT EXISTS likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			post_id INTEGER NOT NULL REFERENCES posts(id),
			UNIQUE(username, post_id)
		);`,

		// Disikes
		`CREATE TABLE IF NOT EXISTS dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			post_id INTEGER NOT NULL REFERENCES posts(id),
			UNIQUE(username, post_id)
		);`,

		// Sessions
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
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
		Db.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, c)
	}
}
