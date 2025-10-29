package comments

import (
	"database/sql"
	"fmt"
	"forum/backend/database"
)

func insertComment(postID, username, content string) error {
		Db:=databasecreate.Open()
	_, err := Db.Exec(
		`INSERT INTO comments (post_id, username, content) VALUES (?, ?, ?)`,
		postID, username, content,
	)
	if err != nil {
		return err
	}

	_, err = Db.Exec(
		`UPDATE posts SET comments_count = comments_count + 1 WHERE id = ?`,
		postID,
	)
	return err
}

func insertCommentLike(commentID int, username string) error {
	Db:=databasecreate.Open()
	var exists int
	err := Db.QueryRow(
		`SELECT 1 FROM comment_likes WHERE comment_id = ? AND username = ?`,
		commentID, username,
	).Scan(&exists)

	if err == nil {
		return fmt.Errorf("already liked")
	}
	if err != sql.ErrNoRows {
		return err
	}

	_, err = Db.Exec(`INSERT INTO comment_likes (username, comment_id) VALUES (?, ?)`, username, commentID)
	if err != nil {
		return err
	}

	_, err = Db.Exec(`UPDATE comments SET likes_count = likes_count + 1 WHERE id = ?`, commentID)
	return err
}
