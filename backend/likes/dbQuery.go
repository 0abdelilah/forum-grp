package likes

import (
	"database/sql"
	"forum/backend/database"
)

// insert like using postid, username
func insertLike(postID, username string) error {
	Db:=databasecreate.Open()
	var exists int
	err := Db.QueryRow(
		`SELECT 1 FROM likes WHERE post_id = ? AND username = ?`,
		postID, username,
	).Scan(&exists)

	if err == nil {
		// user already liked = remove like
		_, err = Db.Exec(
			`DELETE FROM likes WHERE post_id = ? AND username = ?`,
			postID, username,
		)
		if err != nil {
			return err
		}
		_, err = Db.Exec(
			`UPDATE posts SET likes_count = likes_count - 1 WHERE id = ?`,
			postID,
		)
		return err
	}
	if err != sql.ErrNoRows {
		return err // real DB error
	}

	// insert like
	_, err = Db.Exec(
		`INSERT INTO likes (username, post_id) VALUES (?, ?)`,
		username, postID,
	)
	if err != nil {
		return err
	}

	// update counter
	_, err = Db.Exec(
		`UPDATE posts SET likes_count = likes_count + 1 WHERE id = ?`,
		postID,
	)
	return err
}

func insertDislike(postID, username string) error {
	Db:=databasecreate.Open()
	var exists int
	err := Db.QueryRow(
		`SELECT 1 FROM dislikes WHERE post_id = ? AND username = ?`,
		postID, username,
	).Scan(&exists)

	if err == nil {
		// user already liked = remove like
		_, err = Db.Exec(
			`DELETE FROM dislikes WHERE post_id = ? AND username = ?`,
			postID, username,
		)
		if err != nil {
			return err
		}
		_, err = Db.Exec(
			`UPDATE posts SET dislikes_count = dislikes_count - 1 WHERE id = ?`,
			postID,
		)
		return err
	}
	if err != sql.ErrNoRows {
		return err // real DB error
	}

	// insert like
	_, err = Db.Exec(
		`INSERT INTO dislikes (username, post_id) VALUES (?, ?)`,
		username, postID,
	)
	if err != nil {
		return err
	}

	// update counter
	_, err = Db.Exec(
		`UPDATE posts SET dislikes_count = dislikes_count + 1 WHERE id = ?`,
		postID,
	)
	return err
}
