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
		log.Fatal("Failed to open database:", err)
	}

	createTables()
	insertFortrain()
}

// createTables runs all CREATE TABLE statements.
func createTables() {
	createUsersTable()
	createCategoriesTable()
	createPostsTable()
	createCommentsTable()
	createLikesTable()
	createSessionsTable()
}

func createUsersTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);`
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatal("Error creating users table:", err)
	}
}

func createCategoriesTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		categories TEXT NOT NULL,
		post_id INTEGER REFERENCES posts(id)
	);`
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatal("Error creating categories table:", err)
	}

	// Insert default categories (ignore duplicates)
	defaults := []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"}
	for _, v := range defaults {
		_, _ = Db.Exec(`INSERT OR IGNORE INTO categories (categories) VALUES (?)`, v)
	}
}

func createPostsTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL CHECK(length(title) > 0),
		content TEXT NOT NULL CHECK(length(content) > 0),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		comments_count INTEGER DEFAULT 0,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatal("Error creating posts table:", err)
	}
}

func createCommentsTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatal("Error creating comments table:", err)
	}
}

func createLikesTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL CHECK(length(title) > 0),
		content TEXT NOT NULL CHECK(length(content) > 0)
	);`
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatal("Error creating likes table:", err)
	}
}

func createSessionsTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		CHECK(expires_at > created_at)
	);`
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatal("Error creating sessions table:", err)
	}
}

func insertFortrain() {
	_, err := Db.Exec(`
	INSERT INTO posts (user_id, title, content, created_at, likes_count, comments_count)
	VALUES (?, ?, ?, ?, ?, ?)
	`, 1, "dev", "The Fourth of July is the United States' celebration of independence.", "createdAt", 2, 9)
	if err != nil {
		log.Println("Skipping training insert:", err)
	}
}
