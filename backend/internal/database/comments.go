package database

type Comment struct {
	PostId  string `json:"postid"`
	UserId  string `json:"userid"`
	Content string `json:"content"`
}

func GetComments(PostId string) ([]Comment, error) {
	rows, err := Db.Query("SELECT user_id, content, created_at FROM comments WHERE post_id = ?", PostId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var user_id, content, created_at string
		if err := rows.Scan(&user_id, &content, &created_at); err != nil {
			return nil, err
		}
		comments = append(comments, Comment{UserId: user_id, Content: content})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func SaveComment(PostId, UserId, content string) error {
	_, err := Db.Exec(
		"INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", PostId, UserId, content,
	)

	return err
}
