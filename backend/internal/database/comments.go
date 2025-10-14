package database

import "log"

type Comment struct {
	PostId  string `json:"postid"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

func CreateCommentsTable() {
	stmt := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        postid TEXT NOT NULL,
        author TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
}

func GetComments(postid, author, comment string) ([]Comment, error) {
	// Stmt := "SELECT comments FROM posts where postid = ?"

	return nil, nil
}

func SaveComment(postid, author, comment string) error {
	Stmt := "INSERT comments FROM posts where postid = ?"

	return nil
}
