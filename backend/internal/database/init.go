package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func Init() (*sql.DB, error) {
	var err error
	Db, err = sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		return nil, err
	}

	if err := CreateUsersTable(); err != nil {
		return nil, err
	}
	if err := CreatePostsTable(); err != nil {
		return nil, err
	}

	if err := CreateCommentsTable(); err != nil {
		return nil, err
	}

	if err := CreateLikesTable(); err != nil {
		return nil, err
	}
	if err := CreateSessionsTablee(); err != nil {
		return nil, err
	}

	return Db, nil
}

func execSQL(sqlStmt string) error {
	_, err := Db.Exec(sqlStmt)
	return err
}

// ---------- Tables ----------

func CreateUsersTable() error {
	stmt := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);`
	return execSQL(stmt)
}

func CreatePostsTable() error {
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
	return execSQL(stmt)
}

func CreateCommentsTable() error {
	stmt := `
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
    content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	return execSQL(stmt)
}

func CreateLikesTable() error {
	stmt := `
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
	title TEXT NOT NULL CHECK(length(title) > 0),
    content TEXT NOT NULL CHECK(length(content) > 0)
);`
	return execSQL(stmt)
}
func CreateSessionsTablee() error {
	stmt := `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK(expires_at > created_at)
);`
	return execSQL(stmt)
}
