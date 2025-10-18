package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/internal/models"
)

func ALlPageData(r *http.Request, handel string) models.PageData {
	Posts := GetAllPosts()
	return models.PageData{AllPosts: Posts}
}

func GetAllPosts() []models.Post {
	db, _ := sql.Open("sqlite3", "./database/sqlite.db")
	rows, _ := db.Query(`
        SELECT id, user_id, title, content, created_at, 
               likes_count, dislikes_count, comments_count
        FROM posts
        ORDER BY created_at DESC
		
    `)
	fmt.Println("ROWS")
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
