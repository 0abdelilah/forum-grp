package database

import (
	"log"
	"time"

	"forum/backend/models"
)

func GetAlllike(val int, username string) []models.Post {
	var posts []models.Post

	rows, err := Db.Query(`
		SELECT p.id, p.title, p.username, p.content, 
		       p.created_at, p.likes_count, p.dislikes_count, p.comments_count
		FROM posts AS p
		INNER JOIN likes AS l ON p.id = l.post_id
		WHERE l.username = ? AND l.value = ?
		ORDER BY l.created_at ASC
	`, username, val)
	if err != nil {
		log.Printf("failed to query liked posts: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.Id, &p.Title, &p.Username, &p.Content,
			&p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentsNum); err != nil {
			log.Printf("error scanning joined row: %v", err)
			continue
		}

		p.Categories, err = getPostCategories(p.Id)
		if err != nil {
			log.Printf("error getting categories for post %d: %v", p.Id, err)
		}

		if t, err := time.Parse(time.RFC3339, p.CreatedAt); err == nil {
			p.CreatedAt = t.Format("02 Jan 2006 15:04:05")
		}

		posts = append(posts, p)
	}

	return posts
}
