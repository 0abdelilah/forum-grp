package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"forum/backend/models"
)

func AllPageData(r *http.Request, handel string) models.PageData {
	poststr := r.URL.Query().Get("postid")
	postId, _ := strconv.Atoi(poststr)
	if handel == "HomeData" {
		Posts := GetAllPosts(poststr)
		Allcate := GetAllcategores()
		return models.PageData{AllPosts: Posts, Categories: Allcate}
	}
	if handel == "postContent" {
		Postcontent := GetPostdetails(postId)
		// Postcategory:=GetAllcategores()
		// comments:=GetAllcomment(1) //not My task

		// fmt.Println(Post)
		return models.PageData{Postcontent: Postcontent}
	}
	return models.PageData{}
}

func GetAllPosts(postId string) []models.Post {
	var Posts []models.Post
	// if len(postId)> 0 {
	// postId, _ := strconv.Atoi(postId)
	// 	Posts = append(Posts, GetonePost(postId))
	// }
	rows, _ := Db.Query(`
        SELECT id, user_id, title, content, created_at, 
               likes_count, dislikes_count, comments_count
        FROM posts
        ORDER BY created_at DESC
    `)

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

func GetonePost(postId int) models.Post {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	row := Db.QueryRow(`
    SELECT id, user_id, title, content, created_at, 
           likes_count, dislikes_count, comments_count
    FROM posts
    WHERE id = ?
`, postId)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	var post models.Post
	err = row.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.CreatedAt,
		&post.Likes,
		&post.Dislikes,
		&post.CommentsNum,
	)
	return post
}

func GetPostdetails(postId int) models.Postcontent {
	db, err := sql.Open("sqlite3", "./database/sqlite.db")
	if err != nil {
		log.Panic("Failed to open DB:", err)
	}
	defer db.Close()

	row := Db.QueryRow(`
    SELECT id, user_id, title, content, created_at, 
           likes_count, dislikes_count, comments_count
    FROM posts
    WHERE id = ?
`, postId)

	if err != nil {
		panic(err)
	}
	var post models.Postcontent
	err = row.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.CreatedAt,
		&post.Likes,
		&post.Dislikes,
		&post.CommentsNum,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No post found with id=2")
		} else {
			log.Panic("Scan error:", err)
		}
	}

	post.Comments = []models.Comment{}
	fmt.Println("This is the post content:", post)
	return post
}

func GetAlllike() {
}

func GetAllcomment(postId int) {
}

func GetAllcategores() []models.Category {
	// db, _ := sql.Open("sqlite3", "./database/sqlite.db")
	var Allcate []models.Category
	// if err!=nil{
	// 	log.Fatal(err)
	// }
	var cate models.Category
	data, _ := Db.Query(`SELECT id, categories FROM categories`)
	// if !data.Next(){
	// 	log.Fatal("noo selection")
	// }
	for data.Next() {
		data.Scan(&cate.Id, &cate.Category)
		Allcate = append(Allcate, cate)
	}
	fmt.Println(Allcate)
	return Allcate
}
