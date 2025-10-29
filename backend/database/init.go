package databasecreate

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Open() *sql.DB {
	Db, err := sql.Open("sqlite3", "./backend/database/sqlite.db")
	if err := Db.Ping(); err != nil {
		log.Fatal("DB open error:", err)
	}
	if err != nil {
		log.Fatal("DB open error:", err)
	}
	defer Db.Close()
	return Db
}

func Init() {
	createUser()
	CreatePosts()
	CreateCate()
	CreateComments()
	Createlikes()
	CreatePostCategories()
	Creatsessions()
}

func ERR(err error) {
	if err != nil {
		log.Fatal("error :", err)
	}
}

func createUser() {
	Db := Open()
	defer Db.Close()
	_, err := Db.Exec(`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			password_hash TEXT NOT NULL
		);`)
	ERR(err)
}

func CreatePosts() {
	Db := Open()
	defer Db.Close()
	_, err := Db.Exec(`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			title TEXT NOT NULL CHECK(length(title) > 0),
			content TEXT NOT NULL CHECK(length(content) > 0),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			likes_count INTEGER DEFAULT 0,
			dislikes_count INTEGER DEFAULT 0,
			comments_count INTEGER DEFAULT 0
		);`)
	ERR(err)
}

func CreateCate() {
	Db := Open()
	defer Db.Close()
	_, err := Db.Exec(`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		);`,
	)
	InsertCat()
	ERR(err)
}

func CreateComments() {
	Db := Open()
	defer Db.Close()
	_, err := Db.Exec(
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	)
	ERR(err)
}

func CreatePostCategories() {
	Db := Open()
	defer Db.Close()
	_, err := Db.Exec(`CREATE TABLE IF NOT EXISTS post_categories (
			post_id INTEGER NOT NULL REFERENCES posts(id),
			category_id INTEGER NOT NULL REFERENCES categories(id),
			UNIQUE(post_id, category_id)
		);`,
	)
	ERR(err)
}

func Creatsessions() {
	Db := Open()
	defer Db.Close()
	_, err := Db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			expires_at DATETIME NOT NULL CHECK(expires_at > created_at),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	)
	ERR(err)
}

func Createlikes() {
}

func InsertCat() {
	Db := Open()
	defer Db.Close()
	defaults := []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"}
	for _, c := range defaults {
		Db.Exec(`INSERT OR IGNORE INTO categories (name) VALUES (?)`, c)
	}
}
