package likes

import (
	"database/sql"
	"fmt"
	"forum/backend/database"
)

// insert like using postid, username

func insertLike(postID, username string) error {
	var exists int
	err := database.Db.QueryRow(
		`SELECT 1 FROM likes WHERE post_id = ? AND username = ?`,
		postID, username,
	).Scan(&exists)

	if err == nil {
		return fmt.Errorf("user already liked")
	}
	if err != sql.ErrNoRows {
		return err // real DB error
	}

	// insert like
	_, err = database.Db.Exec(
		`INSERT INTO likes (username, post_id) VALUES (?, ?)`,
		username, postID,
	)
	if err != nil {
		return err
	}

	// update counter
	_, err = database.Db.Exec(
		`UPDATE posts SET likes_count = likes_count + 1 WHERE id = ?`,
		postID,
	)
	return err
}

func insertDislike(postID, username string) error {
	var exists int
	err := database.Db.QueryRow(
		`SELECT 1 FROM dislikes WHERE post_id = ? AND username = ?`,
		postID, username,
	).Scan(&exists)

	if err == nil {
		return fmt.Errorf("user already liked")
	}
	if err != sql.ErrNoRows {
		return err // real DB error
	}

	// insert like
	_, err = database.Db.Exec(
		`INSERT INTO dislikes (username, post_id) VALUES (?, ?)`,
		username, postID,
	)
	if err != nil {
		return err
	}

	// update counter
	_, err = database.Db.Exec(
		`UPDATE posts SET dislikes_count = dislikes_count + 1 WHERE id = ?`,
		postID,
	)
	return err
}
