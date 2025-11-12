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
	if err := setupDatabase(); err != nil {
		log.Fatalf("Database setup failed: %v", err)
	}
}

// openDatabase opens the SQLite database connection.
func openDatabase() {
	var err error
	Db, err = sql.Open("sqlite3", "./backend/database/sqlite.db")
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal("DB connection failed:", err)
	}
	log.Println(" Database connection established successfully.")
}

// setupDatabase creates all tables and inserts defaults in one transaction.



func setupDatabase() error {
	tx, err := Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Fatalf("panic during setup: %v", p)
		}
	}()

	// Enable foreign keys
	if _, err := tx.Exec("PRAGMA foreign_keys = ON"); err != nil {
		tx.Rollback()
		return err
	}

	// Create tables
	createTable := func(query, name string) error {
		if _, err := tx.Exec(query); err != nil {
			log.Printf(" Failed to create %s table: %v", name, err)
			return err
		}
		log.Printf("Table '%s' ensured.", name)
		return nil
	}

	// tables
	if err := createTable(usersTable, "users"); err != nil {
		tx.Rollback()
		return err
	}
	if err := createTable(postsTable, "posts"); err != nil {
		tx.Rollback()
		return err
	}
	if err := createTable(categoriesTable, "categories"); err != nil {
		tx.Rollback()
		return err
	}
	if err := createTable(commentsTable, "comments"); err != nil {
		tx.Rollback()
		return err
	}
	if err := createTable(postCategoriesTable, "post_categories"); err != nil {
		tx.Rollback()
		return err
	}
	if err := createTable(likesTable, "likes"); err != nil {
		tx.Rollback()
		return err
	}
	if err := createTable(sessionsTable, "sessions"); err != nil {
		tx.Rollback()
		return err
	}

	// insert default categories
	for _, c := range []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"} {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, c); err != nil {
			log.Printf(" Failed to insert category '%s': %v", c, err)
		}
	}
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_likes_user_value ON likes(username, value);`,
		`CREATE INDEX IF NOT EXISTS idx_likes_created_at ON likes(created_at);`,
	}

	for _, q := range indexes {
		if _, err := tx.Exec(q); err != nil {
			log.Printf(" Failed to create index: %v", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Database setup completed successfully.")
	return nil
}

// === TABLE DEFINITIONS ===
var (
	usersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL CHECK(length(username) BETWEEN 4 AND 24),
		email TEXT UNIQUE NOT NULL CHECK(length(email) <= 100),
		password_hash TEXT NOT NULL
	);`

	postsTable = `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
	    title TEXT NOT NULL CHECK(trim(title) != '' AND length(title) <= 30),
	    content TEXT NOT NULL CHECK(trim(content) != '' AND length(content) <= 300),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		comments_count INTEGER DEFAULT 0
	);`

	categoriesTable = `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);`

	commentsTable = `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		content TEXT NOT NULL CHECK(trim(content) != '' AND length(content) <= 300),
		likes_count INTEGER DEFAULT 0,
		dislikes_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`

	postCategoriesTable = `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL REFERENCES posts(id),
		category_id INTEGER NOT NULL REFERENCES categories(id),
		UNIQUE(post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`

	likesTable = `
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
	);`

	sessionsTable = `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		CHECK(expires_at > created_at)
	);`
)
