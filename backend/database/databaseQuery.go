package database

import (
	"database/sql"
	"forum/backend/models"
	"log"
	"net/http"
	"strconv"
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

func GetAlllike(db *sql.DB, target string, userID int) ([]models.Likes, error) {

	rows, err := db.Query(`
        SELECT id, user_id, target_id, value
        FROM likes
        WHERE user_id = ?
		WHERE target_type = ?
        ORDER BY created_at ASC
    `, userID, target)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []models.Likes
	for rows.Next() {
		var l models.Likes
		if err := rows.Scan(&l.UserId, &l.Target_id, &l.Target_type, &l.Value); err != nil {
			return nil, err

		}
		likes = append(likes, l)

	}
	return likes, nil

}
func GetCommentsByuserID(db *sql.DB, userID int) ([]models.Comment, error) {

	rows, err := db.Query(`
        SELECT id, post_id, post_id, content, created_at
        FROM comments
        WHERE user_id = ?
        ORDER BY created_at ASC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.Id, &c.PostId, &c.UserId, &c.Content, &c.Created); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]models.Comment, error) {
	rows, err := db.Query(`
        SELECT id, post_id, user_id, content, created_at
        FROM comments
        WHERE post_id = ?
        ORDER BY created_at ASC
    `, postID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.Id, &c.PostId, &c.UserId, &c.Content, &c.Created); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
