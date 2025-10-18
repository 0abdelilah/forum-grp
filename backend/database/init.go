package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB
func Init() {
	CreateDB()
	CreateCategoriesTable()
	CreateUsersTable()
	CreatePostsTable()
	CreateCommentsTable()
	CreateLikesTable()
	CreateSessionsTablee()
	insertFortrain()
}

func CreateDB() {
	if err := os.MkdirAll("./database", os.ModeSticky|os.ModePerm); err != nil {
		log.Fatal(err)
	}
	os.Create("./database/sqlite.db")
}

// ---------- Tables ----------

func CreateUsersTable() {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);`
	db.Exec(stmt)
	db.Close()
}

func CreateCategoriesTable() {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt := `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		categories TEXT NOT NULL,
		post_id INTEGER REFERENCES posts(id)
	);`

	_, err = db.Exec(stmt)
	for _, v := range []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"} {
		_, err := db.Exec("INSERT INTO categories (categories) VALUES (?)", v)
		if err != nil {
			log.Fatal("Error inserting category:", err)
		}
	}
	if err != nil {
		log.Fatal("Problem on Category", err)
	}
	db.Close()
}

func CreatePostsTable() {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
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
	db.Exec(stmt)
	db.Close()
}

func CreateCommentsTable() {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt := `
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
    content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	db.Exec(stmt)
	db.Close()
}

func CreateLikesTable() {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt := `
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
	title TEXT NOT NULL CHECK(length(title) > 0),
    content TEXT NOT NULL CHECK(length(content) > 0)
);`
	db.Exec(stmt)
	db.Close()
}

func CreateSessionsTablee() {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	stmt := `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK(expires_at > created_at)
);`
	db.Exec(stmt)
	db.Close()
}

func insertFortrain() {
	db, _ := sql.Open("sqlite3", "./database/sqlite.db")
	_, err := db.Exec(`
    INSERT INTO posts (user_id, title, content, created_at, likes_count, comments_count)
     VALUES (?, ?, ?, ?, ?, ?)
 `, 1, "dev","the The Fourth of July is the United States' celebration of its independence from England.", "createdAt", 2, 9)
	if err != nil {
		log.Fatal("<oschkil finsert:", err)
	}
}
