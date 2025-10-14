package database

import (
	"database/sql"

	"forum/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")

	database.CreateCommentsTable()

	return db, err
}
