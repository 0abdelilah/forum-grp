package database

import "forum/backend/models"

func getComments(postID int) ([]models.Comment, error) {
	rows, err := Db.Query(
		`SELECT id, username, content, likes_count, dislikes_count, created_at FROM comments WHERE post_id = ? ORDER BY created_at DESC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.Id, &c.Username, &c.Content, &c.Likes, &c.Dislikes, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.CreatedAt = PrettifyCreatedAt(c.CreatedAt)
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
