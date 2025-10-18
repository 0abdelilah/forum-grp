package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "forum/internal/auth"
	"forum/internal/comments"
	"forum/internal/database"
	"forum/internal/home"
	"forum/internal/posts"
)

func main() {
	db, err := database.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home.HomeHandler)

	// Authentication
	mux.HandleFunc("GET /register/", handlers.RegisterHandlerGet)
	mux.HandleFunc("POST /register/", handlers.RegisterHandlerPost)

	mux.HandleFunc("GET  /login/", handlers.LoginHandlerGet)
	mux.HandleFunc("POST  /login/", handlers.LoginHandlerPost)

	// Comments
	mux.HandleFunc("POST /api/comment", comments.SaveCommentHandler)
	mux.HandleFunc("GET /api/comments", comments.GetCommentsHandler)

	// Posts
	mux.HandleFunc("GET /api/posts", posts.LoadPostsHandler)
	//mux.HandleFunc("POST /api/create_post", posts.CreatePostsHandler)

	//Posts content 
	mux.HandleFunc("GET /post-detail", posts.SeePostdetail)
	// static files
	mux.HandleFunc("GET /static/", home.StaticHandler)

	fmt.Println("Listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Println(err)
	}
}
