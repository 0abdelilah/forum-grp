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
		log.Fatal("Failed to open database:", err)
	}

	 //CreateTables() //ila bghiti Tables from scrach
	 //InsertInTables() // ila bghit tansirti 
}

// createTables runs all CREATE TABLE statements.
func CreateTables() {
	CreateUsersTable()
	CreateCategoriesTable()
	CreatePostsTable()
	CreateCommentsTable()
	CreateLikesTable()
	CreateSessionsTable()
}
func InsertInTables(){
	InsertPost()
	// InsetComment()
}
func CreateUsersTable() {
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

func CreateCategoriesTable() {
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

func CreatePostsTable() {
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

func CreateCommentsTable() {
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

func CreateLikesTable() {
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

func CreateSessionsTable() {
	stmt := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		CHECK(expires_at > created_at)
	);`
	if _, err := Db.Exec(stmt); err != nil {
		log.Fatal("Error creating sessions table:", err)
	}
}

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
	// db, _ := sql.Open("sqlite3", "./database/sqlite.db")
	_, err := Db.Exec(`
	INSERT INTO comments ( id ,post_id , user_id ,content ,created_at) VALUES(?,?,?,?,?)
	`, 2, 3, 5, "I love you",time.Now().Format("2001-12-3"))
	if err != nil {
		log.Fatal("error in isert:", err)
	}
}
