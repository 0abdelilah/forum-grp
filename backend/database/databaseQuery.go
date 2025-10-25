package database

import (
	"database/sql"
	"fmt"
	"forum/backend/models"
	"log"
	"net/http"
	"strconv"
	"time"
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

	case "postContent":
		post := GetPostDetails(postId)
		return models.PageData{PostContent: post}

	default:
		return models.PageData{}
	}
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
	comments, err := getComments(postId)
	if err != nil {
		fmt.Println("Error getting comments")
	} else {
		post.Comments = comments
	}

	return post
}

func getComments(postID int) ([]models.Comment, error) {
	rows, err := Db.Query(
		`SELECT id, username, content, created_at FROM comments WHERE post_id = ? ORDER BY created_at DESC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.Id, &c.Username, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
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

func GetAllliketarget(db *sql.DB, Target_id int) ([]models.LikesID, error) {

	rows, err := db.Query(`
        SELECT user_id, target_id, value,id
        FROM likes
      WHERE target_id = ?
        ORDER BY created_at ASC
    `, Target_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []models.LikesID
	for rows.Next() {
		var l models.LikesID
		if err := rows.Scan(&l.UserId, &l.Target_id, &l.Value, &l.Id); err != nil {
			return nil, err

		}
		likes = append(likes, l)

	}
	return likes, nil

}
