package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/backend/auth"
	"forum/backend/database"
	"forum/backend/home"
	"forum/backend/likes"
	"forum/backend/posts"
)

func main() {
	database.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home.HomeHandler)
	mux.HandleFunc("POST /", home.PostHomeHandler)

	// Authentication
	mux.HandleFunc("GET /register", auth.RegisterHandlerGet)
	mux.HandleFunc("POST /api/register", auth.RegisterHandlerPost)

	mux.HandleFunc("GET  /login", auth.LoginHandlerGet)
	mux.HandleFunc("POST  /api/login", auth.LoginHandlerPost)

	mux.HandleFunc("GET  /logout", auth.LogoutHandler)

	// Posts
	mux.HandleFunc("POST /api/create_post", posts.CreatePostsHandler)

	mux.HandleFunc("POST /api/like", likes.LikeHandler)
	mux.HandleFunc("POST /api/dislike", likes.DislikeHandler)

	// Posts content
	mux.HandleFunc("GET /post-detail/", posts.SeePostdetail)

	// static files
	mux.HandleFunc("GET /static/", home.StaticHandler)

	fmt.Println("Listening on http://localhost:3001")
	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Println(err)
	}
}
