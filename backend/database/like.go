package database

import (
	"forum/backend/models"
	"log"
	"time"
)

func GetAlllike(val int, target string, username string) []models.LikedPosts {
	var posts []models.LikedPosts

	rows, err := Db.Query(`
        SELECT target_id
        FROM likes
        WHERE username = ? AND target_type = ? AND value = ?
        ORDER BY created_at ASC
    `, username, target, val)
	if err != nil {
		log.Printf("failed to query likes: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var targetID int
		if err := rows.Scan(&targetID); err != nil {
			log.Printf("error scanning like: %v", err)
			continue
		}

		likedPosts := GetAllLikedPosts(targetID)
		posts = append(posts, likedPosts...)
	}

	return posts
}

func GetAllLikedPosts(id int)[] models.LikedPosts {
	var posts []models.LikedPosts

	rows, err := Db.Query(`
		SELECT id, title, username, content, created_at,
		       likes_count, dislikes_count, comments_count
		FROM posts
		WHERE id = ?
		ORDER BY created_at DESC
	`, id)
	if err != nil {
		log.Printf("failed to query posts: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p models.LikedPosts
		if err := rows.Scan(&p.Id, &p.Title, &p.Username, &p.Content,
			&p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentsNum); err != nil {
			log.Printf("error scanning post row: %v", err)
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
