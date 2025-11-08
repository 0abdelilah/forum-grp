package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// Init initializes the database connection and creates all tables and defaults.
func Init() {
	openDatabase()
	createAllTables()
	insertDefaultCategories()
}

// openDatabase opens the SQLite database connection.
func openDatabase() {
	var err error
	Db, err = sql.Open("sqlite3", "./backend/database/sqlite.db")
	if err != nil {
		log.Fatal("DB open error:", err)
	}
}

// createAllTables creates all necessary tables if they don't exist.
func createAllTables() {
	Db.Exec("PRAGMA foreign_keys = ON")

	createUsersTable()
	createPostsTable()
	createCategoriesTable()
	createCommentsTable()
	createPostCategoriesTable()
	createLikesTable()
	createSessionsTable()
}

// createUsersTable creates the users table.
func createUsersTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		email TEXT UNIQUE,
		password_hash TEXT NOT NULL
	);`
	execStmt(stmt, "users")
}

// createPostsTable creates the posts table.
func createPostsTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		title TEXT NOT NULL CHECK(length(title) > 0),
		content TEXT NOT NULL CHECK(length(content) > 0),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		comments_count INTEGER DEFAULT 0
	);`
	execStmt(stmt, "posts")
}

// createCategoriesTable creates the categories table.
func createCategoriesTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`
	execStmt(stmt, "categories")
}

// createCommentsTable creates the comments table.
func createCommentsTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		content TEXT NOT NULL,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	execStmt(stmt, "comments")
}

// createPostCategoriesTable creates the post_categories table.
func createPostCategoriesTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL REFERENCES posts(id),
		category_id INTEGER NOT NULL REFERENCES categories(id),
		UNIQUE(post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	execStmt(stmt, "post_categories")
}

// createLikesTable creates the likes table.
func createLikesTable() {
	stmt := `
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    value INTEGER NOT NULL CHECK(value IN (1, -1)), -- 1 = like, -1 = dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    UNIQUE(username, post_id, comment_id)
);
`
	execStmt(stmt, "likes")
}

// createSessionsTable creates the sessions table.
func createSessionsTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		expires_at DATETIME NOT NULL CHECK(expires_at > created_at),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	execStmt(stmt, "sessions")
}

// insertDefaultCategories seeds default categories.
func insertDefaultCategories() {
	defaults := []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"}
	for _, c := range defaults {
		if _, err := Db.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, c); err != nil {
			log.Println("Failed to insert default category:", c, err)
		}
	}
}

// execStmt executes a SQL statement and logs if it fails.
func execStmt(stmt string, tableName string) {
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatalf("Table creation failed for %s: %v", tableName, err)
	}
}
