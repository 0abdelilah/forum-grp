package database

import (
	"forum/backend/models"
	"log"
	"net/http"
)

func AllPageData(r *http.Request, handel string) models.PageData {
	Posts := GetAllPosts()
	return models.PageData{AllPosts: Posts}
}

func GetAllPosts() []models.Post {
	rows, _ := Db.Query(`
        SELECT id, user_id, title, content, created_at, 
               likes_count, dislikes_count, comments_count
        FROM posts
        ORDER BY created_at DESC
    `)

	var Posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.Id, &p.UserId, &p.Title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes, &p.CommentsNum)
		if err != nil {
			log.Fatal(err)
		}

		Posts = append(Posts, p)
	}

	return Posts
}

func GetPostdetails() {
}

func GetAlllike() {
}

func GetAllcomment() {
}

func GetAllcategores() {
}
