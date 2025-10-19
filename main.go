package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/backend/auth"
	"forum/backend/comments"
	"forum/backend/database"
	"forum/backend/home"
	"forum/backend/posts"
)

func main() {
	//_, err := os.Stat("../frontend/templates/index.html")

	database.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home.HomeHandler)

	// Authentication
	mux.HandleFunc("GET /register/", auth.RegisterHandlerGet)
	mux.HandleFunc("POST /register/", auth.RegisterHandlerPost)

	mux.HandleFunc("GET  /login/", auth.LoginHandlerGet)
	mux.HandleFunc("POST  /login/", auth.LoginHandlerPost)

	// Comments
	mux.HandleFunc("POST /api/comment", comments.SaveCommentHandler)
	mux.HandleFunc("GET /api/comments", comments.GetCommentsHandler)

	// Posts
	mux.HandleFunc("GET /api/posts", posts.LoadPostsHandler)
	// mux.HandleFunc("POST /api/create_post", posts.CreatePostsHandler)

	// Posts content
	mux.HandleFunc("GET /post-detail", posts.SeePostdetail)
	// static files
	mux.HandleFunc("GET /static/", home.StaticHandler)

	fmt.Println("Listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Println(err)
	}
}
