package database

import (
	"database/sql"
	"fmt"
	"forum/backend/models"
	"log"
	"strings"
	"time"
)

func GetPosts(categories []string) []models.Post {
	if categories == nil {
		fmt.Println("Getting all posts")
		return GetAllPosts()
	}
	fmt.Println("Getting filtered posts", categories)
	return GetPostsByCategories(categories)
}

func GetAllPosts() []models.Post {
	var posts []models.Post

	rows, err := Db.Query(`
		SELECT id, title, username, content, created_at, 
			   likes_count, dislikes_count, comments_count
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("failed to query posts: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.Id, &p.Username, &p.Title, &p.Content,
			&p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentsNum); err != nil {
			log.Printf("error scanning post row: %v", err)
			continue
		}

		p.Categories2, err = getPostCategories(p.Id)
		if err != nil {
			fmt.Println("Error getting categories")
		}

		t, _ := time.Parse(time.RFC3339, "2025-10-23T12:34:22Z")
		p.CreatedAt = (t.Format("Mon, Jan 2 2006 • 15:04 MST"))

		posts = append(posts, p)
	}

	return posts
}

func GetPostDetails(postId int) models.Post {
	var post models.Post

	err := Db.QueryRow(`
		SELECT id, username, title, content, created_at,
			   likes_count, dislikes_count, comments_count
		FROM posts
		WHERE id = ?
	`, postId).Scan(&post.Id, &post.Username, &post.Title, &post.Content,
		&post.CreatedAt, &post.Likes, &post.Dislikes, &post.CommentsNum)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no post found with id=%d", postId)
			return models.Post{}
		}
		log.Printf("error fetching post: %v", err)
		return models.Post{}
	}
	post.Comments, err = getComments(postId)
	if err != nil {
		fmt.Println("Error getting comments")
	}

	return post
}

func GetPostsByCategories(categoryNames []string) []models.Post {
	if len(categoryNames) == 0 {
		return nil
	}

	placeholders := strings.Repeat("?,", len(categoryNames))
	placeholders = placeholders[:len(placeholders)-1] // remove last comma

	args := make([]interface{}, len(categoryNames))
	for i, c := range categoryNames {
		args[i] = c
	}

	query := fmt.Sprintf(`
        SELECT DISTINCT p.id, p.title, p.username, p.content, p.created_at,
               p.likes_count, p.dislikes_count, p.comments_count
        FROM posts p
        JOIN post_categories pc ON p.id = pc.post_id
        JOIN categories c ON pc.category_id = c.id
        WHERE c.name IN (%s)
        ORDER BY p.created_at DESC
    `, placeholders)

	rows, err := Db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.Id, &p.Title, &p.Username, &p.Content,
			&p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentsNum); err != nil {
			continue
		}
		p.Categories2, _ = getPostCategories(p.Id)
		posts = append(posts, p)
	}

	return posts
}
