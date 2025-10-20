package database

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"forum/backend/models"
)

func AllPageData(r *http.Request, handle string) models.PageData {
	postStr := r.URL.Query().Get("postid")
	postId, err := strconv.Atoi(postStr)
	if err != nil && postStr != "" {
		log.Printf("invalid postid: %v", err)
		return models.PageData{}
	}

	switch handle {
	case "HomeData":
		posts := GetAllPosts()
		categories := GetAllCategories()
		return models.PageData{AllPosts: posts, Categories: categories}

	case "Post":
		post := GetPostDetails(postId)
		return models.PageData{Postcontent: post}

	default:
		return models.PageData{}
	}
}

func GetAllPosts() []models.Post {
	var posts []models.Post

	rows, err := Db.Query(`
		SELECT id, user_id, title, content, created_at, 
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
		if err := rows.Scan(&p.Id, &p.UserId, &p.Title, &p.Content,
			&p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentsNum); err != nil {
			log.Printf("error scanning post row: %v", err)
			continue
		}
		posts = append(posts, p)
	}
	return posts
}

func GetPostDetails(postId int) models.Post {
	var post models.Post

	err := Db.QueryRow(`
		SELECT id, user_id, title, content, created_at,
			   likes_count, dislikes_count, comments_count
		FROM posts
		WHERE id = ?
	`, postId).Scan(&post.Id, &post.UserId, &post.Title, &post.Content,
		&post.CreatedAt, &post.Likes, &post.Dislikes, &post.CommentsNum)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no post found with id=%d", postId)
			return models.Post{}
		}
		log.Printf("error fetching post: %v", err)
		return models.Post{}
	}

	post.Comments = []models.Comment{}
	return post
}

func GetAllCategories() []models.Category {
	var categories []models.Category

	rows, err := Db.Query(`SELECT id, categories FROM categories`)
	if err != nil {
		log.Printf("failed to query categories: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.Id, &c.Category); err != nil {
			log.Printf("error scanning category: %v", err)
			continue
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
	}
	return categories
}

func GetAlllike() {
}

func GetAllcomment(postId int) {
}
