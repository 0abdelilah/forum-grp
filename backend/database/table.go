// database/schema.go
package database

import (
	"database/sql"
	"fmt"
)

// InitSchema runs all create statements (tables, triggers, indexes, view, seed).
// Call this once at startup (or when you want to initialize the DB).
func InitSchema(db *sql.DB) error {
	// ensure foreign keys ON
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := CreateUsersTable(db); err != nil {
		return err
	}
	if err := CreateCategoriesTable(db); err != nil {
		return err
	}
	if err := CreatePostsTable(db); err != nil {
		return err
	}
	if err := CreatePostCategoriesTable(db); err != nil {
		return err
	}
	if err := CreateCommentsTable(db); err != nil {
		return err
	}
	if err := CreateLikesTable(db); err != nil {
		return err
	}
	if err := CreateSessionsTable(db); err != nil {
		return err
	}

	// indexes
	if err := CreateIndexes(db); err != nil {
		return err
	}
	// seed categories
	if err := SeedCategories(db); err != nil {
		return err
	}

	return nil
}

func execSQL(db *sql.DB, sqlStmt string) error {
	_, err := db.Exec(sqlStmt)
	return err
}

// ---------- Tables ----------

func CreateUsersTable(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	return execSQL(db, stmt)
}

func CreateCategoriesTable(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);`
	return execSQL(db, stmt)
}

func CreatePostsTable(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL CHECK(length(title) > 0),
    content TEXT NOT NULL CHECK(length(content) > 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    likes_count INTEGER DEFAULT 0,
    dislikes_count INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);`
	return execSQL(db, stmt)
}

func CreatePostCategoriesTable(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
);`
	return execSQL(db, stmt)
}

func CreateCommentsTable(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL CHECK(length(content) > 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);`
	return execSQL(db, stmt)
}

func CreateLikesTable(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    target_type TEXT NOT NULL CHECK(target_type IN ('post','comment')),
    target_id INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK(value IN (1, -1)),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, target_type, target_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);`
	return execSQL(db, stmt)
}

func CreateSessionsTable(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK(expires_at > created_at)
);`
	return execSQL(db, stmt)
}

// ---------- Indexes ----------

func CreateIndexes(db *sql.DB) error {
	stmts := []string{
		`CREATE INDEX IF NOT EXISTS idx_posts_user ON posts(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_post_categories_cat ON post_categories(category_id);`,
		`CREATE INDEX IF NOT EXISTS idx_post_categories_post ON post_categories(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_likes_target ON likes(target_type, target_id);`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);`,
	}
	for _, s := range stmts {
		if err := execSQL(db, s); err != nil {
			return fmt.Errorf("create index failed: %w (stmt: %s)", err, s)
		}
	}
	return nil
}

// ---------- Seed ----------

func SeedCategories(db *sql.DB) error {
	stmt := `INSERT OR IGNORE INTO categories (name) VALUES ('General'), ('Programming'), ('Announcements'), ('Off-topic');`
	return execSQL(db, stmt)
}

// ---------- Views ----------
