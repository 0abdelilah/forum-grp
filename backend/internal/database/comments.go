package database

import (
	"fmt"
	"log"
)

type Comment struct {
	PostId  string `json:"postid"`
	UserId  string `json:"userid"`
	Content string `json:"content"`
}

func GetComments(PostId string) ([]Comment, error) {
	rows, err := Db.Query("SELECT * FROM comments WHERE post_id = ?", PostId)

	for rows.Next() {
		var id int
		var author, content string
		if err := rows.Scan(&id, &author, &content); err != nil {
			log.Println(err)
		}
		fmt.Println(id, author, content)
	}

	return nil, err
}

func SaveComment(PostId, UserId, content string) error {
	_, err := Db.Exec(
		"INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", PostId, UserId, content,
	)

	return err
}
